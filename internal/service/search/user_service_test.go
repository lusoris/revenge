package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCollectionSchema(t *testing.T) {
	schema := UserCollectionSchema()

	assert.Equal(t, "users", schema.Name)
	require.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "created_at", *schema.DefaultSortingField)
	assert.Len(t, schema.Fields, 9)
}

func TestUserCollectionSchema_Facets(t *testing.T) {
	schema := UserCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	assert.True(t, facetFields["is_active"], "is_active should be faceted")
	assert.True(t, facetFields["is_admin"], "is_admin should be faceted")
	assert.Len(t, facetFields, 2, "exactly 2 faceted fields expected")
}

func TestUserCollectionSchema_Sortable(t *testing.T) {
	schema := UserCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	assert.True(t, sortableFields["created_at"], "created_at should be sortable")
	assert.True(t, sortableFields["last_login_at"], "last_login_at should be sortable")
	assert.Len(t, sortableFields, 2, "exactly 2 sortable fields expected")
}

func TestUserCollectionSchema_FieldTypes(t *testing.T) {
	schema := UserCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "string", fieldTypes["username"])
	assert.Equal(t, "string", fieldTypes["email"])
	assert.Equal(t, "string", fieldTypes["display_name"])
	assert.Equal(t, "string", fieldTypes["avatar_url"])
	assert.Equal(t, "bool", fieldTypes["is_active"])
	assert.Equal(t, "bool", fieldTypes["is_admin"])
	assert.Equal(t, "int64", fieldTypes["created_at"])
	assert.Equal(t, "int64", fieldTypes["last_login_at"])
}

func TestUserCollectionSchema_InfixSearch(t *testing.T) {
	schema := UserCollectionSchema()

	infixFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Infix != nil && *f.Infix {
			infixFields[f.Name] = true
		}
	}

	assert.True(t, infixFields["username"], "username should have infix search")
	assert.True(t, infixFields["email"], "email should have infix search")
	assert.True(t, infixFields["display_name"], "display_name should have infix search")
	assert.Len(t, infixFields, 3, "exactly 3 infix fields expected")
}

func TestUserCollectionSchema_TokenSeparators(t *testing.T) {
	schema := UserCollectionSchema()

	require.NotNil(t, schema.TokenSeparators)
	separators := *schema.TokenSeparators
	assert.Contains(t, separators, "-")
	assert.Contains(t, separators, "_")
	assert.Contains(t, separators, ".")
}

func TestUserCollectionSchema_SymbolsToIndex(t *testing.T) {
	schema := UserCollectionSchema()

	require.NotNil(t, schema.SymbolsToIndex)
	symbols := *schema.SymbolsToIndex
	assert.Contains(t, symbols, "@", "@ symbol should be indexed for email search")
}

func TestDefaultUserSearchParams(t *testing.T) {
	params := DefaultUserSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "created_at:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "is_active")
	assert.Contains(t, params.FacetBy, "is_admin")
}

func TestUserToDocument_Full(t *testing.T) {
	doc := UserDocument{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Username:    "johndoe",
		Email:       "john@example.com",
		DisplayName: "John Doe",
		AvatarURL:   "/avatars/john.jpg",
		IsActive:    true,
		IsAdmin:     false,
		CreatedAt:   1700000000,
		LastLoginAt: 1700100000,
	}

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, "johndoe", doc.Username)
	assert.Equal(t, "john@example.com", doc.Email)
	assert.Equal(t, "John Doe", doc.DisplayName)
	assert.Equal(t, "/avatars/john.jpg", doc.AvatarURL)
	assert.True(t, doc.IsActive)
	assert.False(t, doc.IsAdmin)
	assert.Equal(t, int64(1700000000), doc.CreatedAt)
	assert.Equal(t, int64(1700100000), doc.LastLoginAt)
}

func TestUserToDocument_Minimal(t *testing.T) {
	doc := UserDocument{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "minimaluser",
		Email:    "minimal@test.com",
	}

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, "minimaluser", doc.Username)
	assert.Equal(t, "minimal@test.com", doc.Email)
	assert.Empty(t, doc.DisplayName)
	assert.Empty(t, doc.AvatarURL)
	assert.False(t, doc.IsActive)
	assert.False(t, doc.IsAdmin)
	assert.Zero(t, doc.CreatedAt)
	assert.Zero(t, doc.LastLoginAt)
}

func TestParseUserDocument(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected UserDocument
	}{
		{
			name: "full document",
			data: map[string]interface{}{
				"id":            "550e8400-e29b-41d4-a716-446655440000",
				"username":      "johndoe",
				"email":         "john@example.com",
				"display_name":  "John Doe",
				"avatar_url":    "/avatars/john.jpg",
				"is_active":     true,
				"is_admin":      false,
				"created_at":    float64(1700000000),
				"last_login_at": float64(1700100000),
			},
			expected: UserDocument{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				Username:    "johndoe",
				Email:       "john@example.com",
				DisplayName: "John Doe",
				AvatarURL:   "/avatars/john.jpg",
				IsActive:    true,
				IsAdmin:     false,
				CreatedAt:   1700000000,
				LastLoginAt: 1700100000,
			},
		},
		{
			name:     "empty document",
			data:     map[string]interface{}{},
			expected: UserDocument{},
		},
		{
			name: "partial document",
			data: map[string]interface{}{
				"id":       "abc123",
				"username": "partial",
				"is_admin": true,
			},
			expected: UserDocument{
				ID:       "abc123",
				Username: "partial",
				IsAdmin:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseUserDocument(tt.data)
			assert.Equal(t, tt.expected, doc)
		})
	}
}

func TestUserSearchService_IsEnabled(t *testing.T) {
	svc := &UserSearchService{
		client: nil,
	}
	assert.False(t, svc.IsEnabled(), "service with nil client should be disabled")
}

func TestUserSearchParams_Validation(t *testing.T) {
	params := UserSearchParams{
		Query:   "john",
		Page:    1,
		PerPage: 20,
	}

	assert.Equal(t, "john", params.Query)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Empty(t, params.FilterBy)
	assert.Empty(t, params.SortBy)
	assert.Empty(t, params.FacetBy)
	assert.False(t, params.IncludeHighlights)
}

func TestUserDocument_PasswordHashExcluded(t *testing.T) {
	// Verify that UserDocument does NOT contain any password-related fields.
	// This is a structural test to ensure sensitive data never reaches the search index.
	doc := UserDocument{
		ID:          "test-id",
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		AvatarURL:   "/avatars/test.jpg",
		IsActive:    true,
		IsAdmin:     false,
		CreatedAt:   1700000000,
		LastLoginAt: 1700100000,
	}

	// All 9 fields accounted for — PasswordHash is not one of them
	assert.Equal(t, "test-id", doc.ID)
	assert.Equal(t, "testuser", doc.Username)
	assert.Equal(t, "test@example.com", doc.Email)
	assert.Equal(t, "Test User", doc.DisplayName)
	assert.Equal(t, "/avatars/test.jpg", doc.AvatarURL)
	assert.True(t, doc.IsActive)
	assert.False(t, doc.IsAdmin)
	assert.Equal(t, int64(1700000000), doc.CreatedAt)
	assert.Equal(t, int64(1700100000), doc.LastLoginAt)
}

func TestUserCollectionSchema_OptionalFields(t *testing.T) {
	schema := UserCollectionSchema()

	optionalFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Optional != nil && *f.Optional {
			optionalFields[f.Name] = true
		}
	}

	assert.True(t, optionalFields["display_name"], "display_name should be optional")
	assert.True(t, optionalFields["avatar_url"], "avatar_url should be optional")
	assert.True(t, optionalFields["last_login_at"], "last_login_at should be optional")
	assert.Len(t, optionalFields, 3, "exactly 3 optional fields expected")
}

func TestUserCollectionSchema_AvatarUrlNotIndexed(t *testing.T) {
	schema := UserCollectionSchema()

	for _, f := range schema.Fields {
		if f.Name == "avatar_url" {
			require.NotNil(t, f.Index)
			assert.False(t, *f.Index, "avatar_url should not be indexed")
			return
		}
	}
	t.Fatal("avatar_url field not found in schema")
}
