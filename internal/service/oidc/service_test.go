// Package oidc provides OIDC/SSO authentication services for Revenge Go.
package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

// mockOIDCProviderRepository is a mock implementation for testing.
type mockOIDCProviderRepository struct {
	providers map[uuid.UUID]*domain.OIDCProvider
}

func newMockOIDCProviderRepository() *mockOIDCProviderRepository {
	return &mockOIDCProviderRepository{
		providers: make(map[uuid.UUID]*domain.OIDCProvider),
	}
}

func (m *mockOIDCProviderRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.OIDCProvider, error) {
	p, ok := m.providers[id]
	if !ok {
		return nil, domain.ErrOIDCProviderNotFound
	}
	return p, nil
}

func (m *mockOIDCProviderRepository) GetByName(_ context.Context, name string) (*domain.OIDCProvider, error) {
	for _, p := range m.providers {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, domain.ErrOIDCProviderNotFound
}

func (m *mockOIDCProviderRepository) List(_ context.Context) ([]*domain.OIDCProvider, error) {
	result := make([]*domain.OIDCProvider, 0, len(m.providers))
	for _, p := range m.providers {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockOIDCProviderRepository) ListEnabled(_ context.Context) ([]*domain.OIDCProvider, error) {
	result := make([]*domain.OIDCProvider, 0)
	for _, p := range m.providers {
		if p.Enabled {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockOIDCProviderRepository) Create(_ context.Context, params domain.CreateOIDCProviderParams) (*domain.OIDCProvider, error) {
	p := &domain.OIDCProvider{
		ID:                    uuid.New(),
		Name:                  params.Name,
		DisplayName:           params.DisplayName,
		IssuerURL:             params.IssuerURL,
		ClientID:              params.ClientID,
		ClientSecretEncrypted: params.ClientSecretEncrypted,
		Scopes:                params.Scopes,
		Enabled:               params.Enabled,
		AutoCreateUsers:       params.AutoCreateUsers,
		DefaultAdmin:          params.DefaultAdmin,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
	m.providers[p.ID] = p
	return p, nil
}

func (m *mockOIDCProviderRepository) Update(_ context.Context, params domain.UpdateOIDCProviderParams) (*domain.OIDCProvider, error) {
	p, ok := m.providers[params.ID]
	if !ok {
		return nil, domain.ErrOIDCProviderNotFound
	}
	if params.DisplayName != nil {
		p.DisplayName = *params.DisplayName
	}
	if params.IssuerURL != nil {
		p.IssuerURL = *params.IssuerURL
	}
	if params.ClientID != nil {
		p.ClientID = *params.ClientID
	}
	if params.Enabled != nil {
		p.Enabled = *params.Enabled
	}
	p.UpdatedAt = time.Now()
	return p, nil
}

func (m *mockOIDCProviderRepository) SetEnabled(_ context.Context, id uuid.UUID, enabled bool) error {
	p, ok := m.providers[id]
	if !ok {
		return domain.ErrOIDCProviderNotFound
	}
	p.Enabled = enabled
	return nil
}

func (m *mockOIDCProviderRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.providers, id)
	return nil
}

func (m *mockOIDCProviderRepository) addProvider(p *domain.OIDCProvider) {
	m.providers[p.ID] = p
}

func TestService_GetProviders(t *testing.T) {
	repo := newMockOIDCProviderRepository()

	// Add some providers
	provider1 := &domain.OIDCProvider{
		ID:          uuid.New(),
		Name:        "keycloak",
		DisplayName: "Keycloak SSO",
		IssuerURL:   "https://keycloak.example.com/realms/revenge",
		ClientID:    "revenge",
		Enabled:     true,
	}
	provider2 := &domain.OIDCProvider{
		ID:          uuid.New(),
		Name:        "disabled",
		DisplayName: "Disabled Provider",
		IssuerURL:   "https://disabled.example.com",
		ClientID:    "test",
		Enabled:     false,
	}
	repo.addProvider(provider1)
	repo.addProvider(provider2)

	svc := &Service{
		providers: repo,
	}

	providers, err := svc.GetProviders(context.Background())
	if err != nil {
		t.Fatalf("GetProviders() error = %v", err)
	}

	if len(providers) != 1 {
		t.Errorf("GetProviders() returned %d providers, want 1 (enabled only)", len(providers))
	}

	if len(providers) > 0 && providers[0].Name != "keycloak" {
		t.Errorf("GetProviders() returned provider %s, want keycloak", providers[0].Name)
	}
}

func TestService_GetAuthorizationURL_ProviderNotFound(t *testing.T) {
	repo := newMockOIDCProviderRepository()

	svc := &Service{
		providers:      repo,
		stateStore:     make(map[string]*authState),
		discoveryCache: make(map[string]*cachedDiscovery),
	}

	_, err := svc.GetAuthorizationURL(context.Background(), uuid.New(), "http://localhost/callback")
	if err == nil {
		t.Fatal("GetAuthorizationURL() should return error for non-existent provider")
	}
}

func TestService_GetAuthorizationURL_ProviderDisabled(t *testing.T) {
	repo := newMockOIDCProviderRepository()
	provider := &domain.OIDCProvider{
		ID:        uuid.New(),
		Name:      "disabled",
		IssuerURL: "https://test.example.com",
		ClientID:  "test",
		Enabled:   false,
	}
	repo.addProvider(provider)

	svc := &Service{
		providers:      repo,
		stateStore:     make(map[string]*authState),
		discoveryCache: make(map[string]*cachedDiscovery),
	}

	_, err := svc.GetAuthorizationURL(context.Background(), provider.ID, "http://localhost/callback")
	if err == nil {
		t.Fatal("GetAuthorizationURL() should return error for disabled provider")
	}
}

func TestService_GetAuthorizationURL_Success(t *testing.T) {
	// Create a mock OIDC discovery server
	discoveryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"issuer": "https://test.example.com",
				"authorization_endpoint": "https://test.example.com/auth",
				"token_endpoint": "https://test.example.com/token",
				"jwks_uri": "https://test.example.com/jwks"
			}`))
		}
	}))
	defer discoveryServer.Close()

	repo := newMockOIDCProviderRepository()
	provider := &domain.OIDCProvider{
		ID:        uuid.New(),
		Name:      "test",
		IssuerURL: discoveryServer.URL,
		ClientID:  "test-client",
		Scopes:    []string{"openid", "profile", "email"},
		Enabled:   true,
	}
	repo.addProvider(provider)

	svc := &Service{
		providers:      repo,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		stateStore:     make(map[string]*authState),
		discoveryCache: make(map[string]*cachedDiscovery),
		stateExpiry:    10 * time.Minute,
	}

	url, err := svc.GetAuthorizationURL(context.Background(), provider.ID, "http://localhost/callback")
	if err != nil {
		t.Fatalf("GetAuthorizationURL() error = %v", err)
	}

	if url == "" {
		t.Error("GetAuthorizationURL() returned empty URL")
	}

	// Verify state was stored
	if len(svc.stateStore) != 1 {
		t.Errorf("stateStore should have 1 entry, got %d", len(svc.stateStore))
	}
}

func TestService_CleanupExpiredStates(t *testing.T) {
	svc := &Service{
		stateStore:  make(map[string]*authState),
		stateExpiry: 100 * time.Millisecond,
	}

	// Add some states
	svc.stateStore["valid"] = &authState{
		ProviderID: uuid.New(),
		CreatedAt:  time.Now(),
	}
	svc.stateStore["expired"] = &authState{
		ProviderID: uuid.New(),
		CreatedAt:  time.Now().Add(-200 * time.Millisecond),
	}

	svc.CleanupExpiredStates()

	if len(svc.stateStore) != 1 {
		t.Errorf("CleanupExpiredStates() should leave 1 state, got %d", len(svc.stateStore))
	}

	if _, ok := svc.stateStore["valid"]; !ok {
		t.Error("CleanupExpiredStates() removed the valid state")
	}
}

func TestGenerateRandomString(t *testing.T) {
	s1, err := generateRandomString(32)
	if err != nil {
		t.Fatalf("generateRandomString() error = %v", err)
	}

	if len(s1) != 32 {
		t.Errorf("generateRandomString(32) returned string of length %d", len(s1))
	}

	s2, err := generateRandomString(32)
	if err != nil {
		t.Fatalf("generateRandomString() error = %v", err)
	}

	if s1 == s2 {
		t.Error("generateRandomString() returned same value twice")
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-12345678901234567890"
	challenge := generateCodeChallenge(verifier)

	if challenge == "" {
		t.Error("generateCodeChallenge() returned empty string")
	}

	// Same verifier should produce same challenge
	challenge2 := generateCodeChallenge(verifier)
	if challenge != challenge2 {
		t.Error("generateCodeChallenge() should be deterministic")
	}

	// Different verifier should produce different challenge
	challenge3 := generateCodeChallenge("different-verifier")
	if challenge == challenge3 {
		t.Error("generateCodeChallenge() should produce different challenges for different verifiers")
	}
}

func TestDiscoveryDocument_Caching(t *testing.T) {
	callCount := 0
	discoveryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"issuer": "https://test.example.com",
			"authorization_endpoint": "https://test.example.com/auth",
			"token_endpoint": "https://test.example.com/token",
			"jwks_uri": "https://test.example.com/jwks"
		}`))
	}))
	defer discoveryServer.Close()

	svc := &Service{
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		discoveryCache: make(map[string]*cachedDiscovery),
	}

	// First call should hit the server
	_, err := svc.fetchDiscovery(context.Background(), discoveryServer.URL)
	if err != nil {
		t.Fatalf("fetchDiscovery() error = %v", err)
	}

	if callCount != 1 {
		t.Errorf("First call should hit server once, got %d calls", callCount)
	}

	// Second call should use cache
	_, err = svc.fetchDiscovery(context.Background(), discoveryServer.URL)
	if err != nil {
		t.Fatalf("fetchDiscovery() error = %v", err)
	}

	if callCount != 1 {
		t.Errorf("Second call should use cache, got %d calls", callCount)
	}
}

// Benchmark tests

func BenchmarkGenerateRandomString(b *testing.B) {
	for b.Loop() {
		_, _ = generateRandomString(32)
	}
}

func BenchmarkGenerateCodeChallenge(b *testing.B) {
	verifier := make([]byte, 64)
	_, _ = rand.Read(verifier)
	v := base64.RawURLEncoding.EncodeToString(verifier)

	for b.Loop() {
		generateCodeChallenge(v)
	}
}
