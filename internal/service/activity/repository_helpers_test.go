package activity

import (
	"encoding/json"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// uuidToPgtype Tests
// ============================================================================

func TestUuidToPgtype(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   uuid.UUID
	}{
		{
			name: "valid UUID v7",
			id:   uuid.Must(uuid.NewV7()),
		},
		{
			name: "nil UUID",
			id:   uuid.Nil,
		},
		{
			name: "specific UUID",
			id:   uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := uuidToPgtype(tt.id)
			assert.True(t, result.Valid)
			assert.Equal(t, [16]byte(tt.id), result.Bytes)
		})
	}
}

// ============================================================================
// dbActivityToEntry Tests
// ============================================================================

func TestDbActivityToEntry_AllFields(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	username := "testuser"
	resourceType := "user"
	userAgent := "TestAgent/1.0"
	errorMsg := "some error"
	now := time.Now()
	success := true
	ip := netip.MustParseAddr("192.168.1.1")

	changes, _ := json.Marshal(map[string]interface{}{"field": "value"})
	metadata, _ := json.Marshal(map[string]interface{}{"key": "data"})

	dbLog := db.ActivityLog{
		ID: uuid.Must(uuid.NewV7()),
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
		Username: &username,
		Action:   ActionUserLogin,
		ResourceType: &resourceType,
		ResourceID: pgtype.UUID{
			Bytes: resourceID,
			Valid: true,
		},
		Changes:      changes,
		Metadata:     metadata,
		IpAddress:    ip,
		UserAgent:    &userAgent,
		Success:      &success,
		ErrorMessage: &errorMsg,
		CreatedAt:    now,
	}

	entry := dbActivityToEntry(dbLog)

	require.NotNil(t, entry)
	assert.Equal(t, dbLog.ID, entry.ID)
	assert.Equal(t, ActionUserLogin, entry.Action)
	assert.Equal(t, now, entry.CreatedAt)
	assert.True(t, entry.Success)

	// Check pointer fields are populated
	require.NotNil(t, entry.UserID)
	assert.Equal(t, userID, *entry.UserID)

	require.NotNil(t, entry.ResourceID)
	assert.Equal(t, resourceID, *entry.ResourceID)

	require.NotNil(t, entry.Username)
	assert.Equal(t, "testuser", *entry.Username)

	require.NotNil(t, entry.ResourceType)
	assert.Equal(t, "user", *entry.ResourceType)

	require.NotNil(t, entry.UserAgent)
	assert.Equal(t, "TestAgent/1.0", *entry.UserAgent)

	require.NotNil(t, entry.ErrorMessage)
	assert.Equal(t, "some error", *entry.ErrorMessage)

	require.NotNil(t, entry.IPAddress)

	require.NotNil(t, entry.Changes)
	assert.Equal(t, "value", entry.Changes["field"])

	require.NotNil(t, entry.Metadata)
	assert.Equal(t, "data", entry.Metadata["key"])
}

func TestDbActivityToEntry_MinimalFields(t *testing.T) {
	t.Parallel()

	now := time.Now()
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "system.startup",
		CreatedAt: now,
		// All optional fields are zero-value / nil
	}

	entry := dbActivityToEntry(dbLog)

	require.NotNil(t, entry)
	assert.Equal(t, dbLog.ID, entry.ID)
	assert.Equal(t, "system.startup", entry.Action)
	assert.Equal(t, now, entry.CreatedAt)
	assert.False(t, entry.Success) // nil Success -> false

	// All optional fields should be nil
	assert.Nil(t, entry.UserID)
	assert.Nil(t, entry.ResourceID)
	assert.Nil(t, entry.Username)
	assert.Nil(t, entry.ResourceType)
	assert.Nil(t, entry.UserAgent)
	assert.Nil(t, entry.ErrorMessage)
	assert.Nil(t, entry.IPAddress)
	assert.Nil(t, entry.Changes)
	assert.Nil(t, entry.Metadata)
}

func TestDbActivityToEntry_InvalidUUIDs(t *testing.T) {
	t.Parallel()

	// UserID and ResourceID with Valid=false should result in nil pointers
	dbLog := db.ActivityLog{
		ID:     uuid.Must(uuid.NewV7()),
		Action: "test",
		UserID: pgtype.UUID{
			Valid: false,
		},
		ResourceID: pgtype.UUID{
			Valid: false,
		},
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	assert.Nil(t, entry.UserID)
	assert.Nil(t, entry.ResourceID)
}

func TestDbActivityToEntry_InvalidJSON(t *testing.T) {
	t.Parallel()

	// Invalid JSON for changes and metadata should result in nil maps
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Changes:   []byte("not valid json{{{"),
		Metadata:  []byte("also not valid"),
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	assert.Nil(t, entry.Changes, "invalid JSON changes should be nil")
	assert.Nil(t, entry.Metadata, "invalid JSON metadata should be nil")
}

func TestDbActivityToEntry_EmptyJSON(t *testing.T) {
	t.Parallel()

	// Empty byte slices for changes and metadata should result in nil maps
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Changes:   nil,
		Metadata:  nil,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	assert.Nil(t, entry.Changes)
	assert.Nil(t, entry.Metadata)
}

func TestDbActivityToEntry_EmptyJSONBytes(t *testing.T) {
	t.Parallel()

	// Zero-length byte slices
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Changes:   []byte{},
		Metadata:  []byte{},
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	assert.Nil(t, entry.Changes)
	assert.Nil(t, entry.Metadata)
}

func TestDbActivityToEntry_SuccessNil(t *testing.T) {
	t.Parallel()

	// Success is nil (not set)
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Success:   nil,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)
	assert.False(t, entry.Success, "nil success should default to false")
}

func TestDbActivityToEntry_SuccessTrue(t *testing.T) {
	t.Parallel()

	success := true
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Success:   &success,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)
	assert.True(t, entry.Success)
}

func TestDbActivityToEntry_SuccessFalse(t *testing.T) {
	t.Parallel()

	success := false
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		Success:   &success,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)
	assert.False(t, entry.Success)
}

func TestDbActivityToEntry_IPv4Address(t *testing.T) {
	t.Parallel()

	ip := netip.MustParseAddr("10.0.0.1")
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		IpAddress: ip,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	require.NotNil(t, entry.IPAddress)
	assert.True(t, entry.IPAddress.Equal(ip.AsSlice()))
}

func TestDbActivityToEntry_IPv6Address(t *testing.T) {
	t.Parallel()

	ip := netip.MustParseAddr("2001:db8::1")
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		IpAddress: ip,
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)

	require.NotNil(t, entry.IPAddress)
}

func TestDbActivityToEntry_ZeroIP(t *testing.T) {
	t.Parallel()

	// Zero-value netip.Addr is not valid
	dbLog := db.ActivityLog{
		ID:        uuid.Must(uuid.NewV7()),
		Action:    "test",
		IpAddress: netip.Addr{},
		CreatedAt: time.Now(),
	}

	entry := dbActivityToEntry(dbLog)
	assert.Nil(t, entry.IPAddress, "zero-value IP should result in nil")
}

// ============================================================================
// NewRepositoryPg Tests
// ============================================================================

func TestNewRepositoryPg(t *testing.T) {
	t.Parallel()

	// Test that the constructor works with nil queries (won't panic on creation)
	repo := NewRepositoryPg(nil)
	require.NotNil(t, repo)
	assert.Nil(t, repo.queries)
}

func TestNewRepositoryPg_WithQueries(t *testing.T) {
	t.Parallel()

	queries := &db.Queries{}
	repo := NewRepositoryPg(queries)
	require.NotNil(t, repo)
	assert.Equal(t, queries, repo.queries)
}
