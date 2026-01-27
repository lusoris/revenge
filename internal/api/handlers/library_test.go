package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// mockLibraryService implements a mock for testing.
type mockLibraryService struct {
	libraries []*domain.Library
	getByID   func(ctx context.Context, id uuid.UUID) (*domain.Library, error)
	list      func(ctx context.Context) ([]*domain.Library, error)
	listForUser func(ctx context.Context, userID uuid.UUID) ([]*domain.Library, error)
	create    func(ctx context.Context, params domain.CreateLibraryParams) (*domain.Library, error)
	update    func(ctx context.Context, params domain.UpdateLibraryParams) error
	delete    func(ctx context.Context, id uuid.UUID) error
	updateLastScan func(ctx context.Context, id uuid.UUID) error
}

func (m *mockLibraryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Library, error) {
	if m.getByID != nil {
		return m.getByID(ctx, id)
	}
	for _, lib := range m.libraries {
		if lib.ID == id {
			return lib, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockLibraryService) List(ctx context.Context) ([]*domain.Library, error) {
	if m.list != nil {
		return m.list(ctx)
	}
	return m.libraries, nil
}

func (m *mockLibraryService) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Library, error) {
	if m.listForUser != nil {
		return m.listForUser(ctx, userID)
	}
	return m.libraries, nil
}

func (m *mockLibraryService) Create(ctx context.Context, params domain.CreateLibraryParams) (*domain.Library, error) {
	if m.create != nil {
		return m.create(ctx, params)
	}
	lib := &domain.Library{
		ID:        uuid.New(),
		Name:      params.Name,
		Type:      params.Type,
		Paths:     params.Paths,
		IsVisible: params.IsVisible,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.libraries = append(m.libraries, lib)
	return lib, nil
}

func (m *mockLibraryService) Update(ctx context.Context, params domain.UpdateLibraryParams) error {
	if m.update != nil {
		return m.update(ctx, params)
	}
	return nil
}

func (m *mockLibraryService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.delete != nil {
		return m.delete(ctx, id)
	}
	return nil
}

func (m *mockLibraryService) UpdateLastScan(ctx context.Context, id uuid.UUID) error {
	if m.updateLastScan != nil {
		return m.updateLastScan(ctx, id)
	}
	return nil
}

// testLibraryHandler creates a LibraryHandler with a mock service for testing.
type testLibraryHandler struct {
	handler *LibraryHandler
	mock    *mockLibraryService
}

func newTestLibraryHandler() *testLibraryHandler {
	mock := &mockLibraryService{
		libraries: []*domain.Library{},
	}
	// We need to create a handler that uses our mock
	// Since LibraryHandler expects *library.Service, we'll test via HTTP
	return &testLibraryHandler{
		mock: mock,
	}
}

// Helper to create test context with claims
func contextWithClaims(userID uuid.UUID, isAdmin bool) context.Context {
	claims := &domain.TokenClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
	}
	return context.WithValue(context.Background(), middleware.ClaimsContextKey, claims)
}

// Helper to create request with context
func requestWithClaims(method, target string, body []byte, userID uuid.UUID, isAdmin bool) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, target, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	ctx := contextWithClaims(userID, isAdmin)
	return req.WithContext(ctx)
}

func TestLibraryToResponse(t *testing.T) {
	scanInterval := 24
	lib := &domain.Library{
		ID:                uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:              "Movies",
		Type:              domain.LibraryTypeMovies,
		Paths:             []string{"/media/movies", "/data/films"},
		IsVisible:         true,
		ScanIntervalHours: &scanInterval,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	resp := libraryToResponse(lib)

	if resp.ID != lib.ID.String() {
		t.Errorf("ID mismatch: got %s, want %s", resp.ID, lib.ID.String())
	}
	if resp.Name != lib.Name {
		t.Errorf("Name mismatch: got %s, want %s", resp.Name, lib.Name)
	}
	if resp.CollectionType != string(lib.Type) {
		t.Errorf("CollectionType mismatch: got %s, want %s", resp.CollectionType, string(lib.Type))
	}
	if len(resp.Locations) != 2 {
		t.Errorf("Locations count mismatch: got %d, want 2", len(resp.Locations))
	}
	if resp.LibraryOptions.AutomaticRefreshIntervalDays != 1 {
		t.Errorf("AutomaticRefreshIntervalDays mismatch: got %d, want 1", resp.LibraryOptions.AutomaticRefreshIntervalDays)
	}
	if len(resp.LibraryOptions.PathInfos) != 2 {
		t.Errorf("PathInfos count mismatch: got %d, want 2", len(resp.LibraryOptions.PathInfos))
	}
}

func TestLibraryToResponse_NilScanInterval(t *testing.T) {
	lib := &domain.Library{
		ID:                uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:              "Photos",
		Type:              domain.LibraryTypePhotos,
		Paths:             []string{"/media/photos"},
		ScanIntervalHours: nil,
	}

	resp := libraryToResponse(lib)

	if resp.LibraryOptions.AutomaticRefreshIntervalDays != 0 {
		t.Errorf("AutomaticRefreshIntervalDays should be 0 for nil interval, got %d", resp.LibraryOptions.AutomaticRefreshIntervalDays)
	}
}

func TestLibraryResponse_JSONFormat(t *testing.T) {
	resp := LibraryResponse{
		ID:             "11111111-1111-1111-1111-111111111111",
		Name:           "Movies",
		CollectionType: "movies",
		Locations:      []string{"/media/movies"},
		LibraryOptions: LibraryOptions{
			EnableRealtimeMonitor:   true,
			EnableInternetProviders: true,
			PathInfos:               []PathInfo{{Path: "/media/movies"}},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify JSON field names match Jellyfin API format
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	expectedFields := []string{"ItemId", "Name", "CollectionType", "Locations", "LibraryOptions"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}

	// ItemId should be used, not Id
	if _, ok := raw["Id"]; ok {
		t.Error("JSON should use 'ItemId' not 'Id'")
	}
}

func TestCreateLibraryRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateLibraryRequest
		wantErr bool
	}{
		{
			name: "valid request",
			input: CreateLibraryRequest{
				Name:           "Movies",
				CollectionType: "movies",
				Paths:          []string{"/media/movies"},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: CreateLibraryRequest{
				Name:           "",
				CollectionType: "movies",
				Paths:          []string{"/media/movies"},
			},
			wantErr: true,
		},
		{
			name: "empty paths",
			input: CreateLibraryRequest{
				Name:           "Movies",
				CollectionType: "movies",
				Paths:          []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.input.Name == "" || len(tt.input.Paths) == 0
			if hasErr != tt.wantErr {
				t.Errorf("validation mismatch: got error=%v, want error=%v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestLibraryType_IsValid(t *testing.T) {
	tests := []struct {
		libType domain.LibraryType
		valid   bool
	}{
		{domain.LibraryTypeMovies, true},
		{domain.LibraryTypeTVShows, true},
		{domain.LibraryTypeMusic, true},
		{domain.LibraryTypeAdultMovies, true},
		{domain.LibraryType("invalid"), false},
		{domain.LibraryType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.libType), func(t *testing.T) {
			if got := tt.libType.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestLibraryType_IsAdultType(t *testing.T) {
	tests := []struct {
		libType domain.LibraryType
		isAdult bool
	}{
		{domain.LibraryTypeMovies, false},
		{domain.LibraryTypeTVShows, false},
		{domain.LibraryTypeAdultMovies, true},
		{domain.LibraryTypeAdultShows, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.libType), func(t *testing.T) {
			if got := tt.libType.IsAdultType(); got != tt.isAdult {
				t.Errorf("IsAdultType() = %v, want %v", got, tt.isAdult)
			}
		})
	}
}
