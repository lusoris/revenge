package search

import (
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// UserCollectionName is the Typesense collection name for users.
const UserCollectionName = "users"

// UserDocument represents a user in the search index for admin user management.
type UserDocument struct {
	ID          string `json:"id"`           // User UUID
	Username    string `json:"username"`     // Username
	Email       string `json:"email"`        // Email address
	DisplayName string `json:"display_name"` // Display name
	AvatarURL   string `json:"avatar_url"`   // Avatar URL
	IsActive    bool   `json:"is_active"`    // Account active status
	IsAdmin     bool   `json:"is_admin"`     // Admin flag
	CreatedAt   int64  `json:"created_at"`   // Unix timestamp
	LastLoginAt int64  `json:"last_login_at"` // Unix timestamp of last login
}

// UserCollectionSchema returns the Typesense collection schema for users.
func UserCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                UserCollectionName,
		DefaultSortingField: ptr.To("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"@"},
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "username", Type: "string", Infix: ptr.To(true)},
			{Name: "email", Type: "string", Infix: ptr.To(true)},
			{Name: "display_name", Type: "string", Infix: ptr.To(true), Optional: ptr.To(true)},
			{Name: "avatar_url", Type: "string", Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "is_active", Type: "bool", Facet: ptr.To(true)},
			{Name: "is_admin", Type: "bool", Facet: ptr.To(true)},
			{Name: "created_at", Type: "int64", Sort: ptr.To(true)},
			{Name: "last_login_at", Type: "int64", Sort: ptr.To(true), Optional: ptr.To(true)},
		},
	}
}
