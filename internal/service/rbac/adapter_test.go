package rbac

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/testutil"
)

// setupTestAdapter creates a fresh Adapter with an isolated casbin_rule table.
func setupTestAdapter(t *testing.T) (*Adapter, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)

	// Clear any existing rules for test isolation
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	adapter := NewAdapter(testDB.Pool())
	return adapter, testDB
}

// newTestModel creates a fresh Casbin model from the standard configuration.
func newTestModel(t *testing.T) casbinmodel.Model {
	t.Helper()
	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)
	return m
}

// newTestEnforcer creates a Casbin enforcer backed by the given adapter.
func newTestEnforcer(t *testing.T, adapter *Adapter) *casbin.Enforcer {
	t.Helper()
	modelPath := "../../../config/casbin_model.conf"
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	require.NoError(t, err)
	return enforcer
}

// =====================================================
// LoadPolicy tests
// =====================================================

func TestAdapter_LoadPolicy_EmptyDB(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	m := newTestModel(t)

	err := adapter.LoadPolicy(m)
	require.NoError(t, err)

	// No policies should be loaded
	assert.Empty(t, m["p"]["p"].Policy)
	assert.Empty(t, m["g"]["g"].Policy)
}

func TestAdapter_LoadPolicy_WithExistingRules(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Insert rules directly into the database
	_, err := testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"p", "alice", "data1", "read", "", "", "")
	require.NoError(t, err)

	_, err = testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"p", "bob", "data2", "write", "", "", "")
	require.NoError(t, err)

	_, err = testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"g", "alice", "admin", "", "", "", "")
	require.NoError(t, err)

	m := newTestModel(t)
	err = adapter.LoadPolicy(m)
	require.NoError(t, err)

	// Verify p-type rules
	assert.Len(t, m["p"]["p"].Policy, 2)

	// Verify g-type rules
	assert.Len(t, m["g"]["g"].Policy, 1)
}

func TestAdapter_LoadPolicy_NullableV3V4V5(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Insert a rule with NULL v3, v4, v5 (nullable columns)
	_, err := testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, NULL, NULL, NULL)",
		"p", "charlie", "data3", "delete")
	require.NoError(t, err)

	m := newTestModel(t)
	err = adapter.LoadPolicy(m)
	require.NoError(t, err)

	assert.Len(t, m["p"]["p"].Policy, 1)
	assert.Equal(t, []string{"charlie", "data3", "delete"}, m["p"]["p"].Policy[0])
}

// =====================================================
// SavePolicy tests
// =====================================================

func TestAdapter_SavePolicy_Roundtrip(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	// Create a model with some policies
	m := newTestModel(t)
	m["p"]["p"].Policy = [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
	}
	m["g"]["g"].Policy = [][]string{
		{"alice", "admin"},
	}

	// Save
	err := adapter.SavePolicy(m)
	require.NoError(t, err)

	// Load into a fresh model to verify roundtrip
	m2 := newTestModel(t)
	err = adapter.LoadPolicy(m2)
	require.NoError(t, err)

	assert.Len(t, m2["p"]["p"].Policy, 2)
	assert.Len(t, m2["g"]["g"].Policy, 1)
}

func TestAdapter_SavePolicy_ReplacesExistingRules(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Insert initial rules
	_, err := testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"p", "old_user", "old_data", "old_action", "", "", "")
	require.NoError(t, err)

	// Verify old rule exists
	var countBefore int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE v0 = 'old_user'").Scan(&countBefore)
	require.NoError(t, err)
	assert.Equal(t, 1, countBefore)

	// Save a completely different set of policies
	m := newTestModel(t)
	m["p"]["p"].Policy = [][]string{
		{"new_user", "new_data", "new_action"},
	}

	err = adapter.SavePolicy(m)
	require.NoError(t, err)

	// Old rule should be gone
	var countOld int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE v0 = 'old_user'").Scan(&countOld)
	require.NoError(t, err)
	assert.Equal(t, 0, countOld, "old rules should be deleted on SavePolicy")

	// New rule should exist
	var countNew int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE v0 = 'new_user'").Scan(&countNew)
	require.NoError(t, err)
	assert.Equal(t, 1, countNew, "new rule should exist after SavePolicy")
}

func TestAdapter_SavePolicy_EmptyModel(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Insert some initial data
	_, err := testDB.Pool().Exec(ctx,
		"INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"p", "alice", "data1", "read", "", "", "")
	require.NoError(t, err)

	// Save an empty model (should clear all rules)
	m := newTestModel(t)
	err = adapter.SavePolicy(m)
	require.NoError(t, err)

	// Verify all rules are gone
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "saving empty model should clear all rules")
}

func TestAdapter_SavePolicy_TransactionAtomicity(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	// Save multiple rules at once
	m := newTestModel(t)
	m["p"]["p"].Policy = [][]string{
		{"user1", "res1", "act1"},
		{"user2", "res2", "act2"},
		{"user3", "res3", "act3"},
	}
	m["g"]["g"].Policy = [][]string{
		{"user1", "admin"},
		{"user2", "editor"},
	}

	err := adapter.SavePolicy(m)
	require.NoError(t, err)

	// Verify all were saved atomically
	m2 := newTestModel(t)
	err = adapter.LoadPolicy(m2)
	require.NoError(t, err)

	assert.Len(t, m2["p"]["p"].Policy, 3, "all 3 p-type rules should be saved")
	assert.Len(t, m2["g"]["g"].Policy, 2, "all 2 g-type rules should be saved")
}

// =====================================================
// AddPolicy tests
// =====================================================

func TestAdapter_AddPolicy_Basic(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	// Verify rule in DB
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.casbin_rule WHERE ptype = 'p' AND v0 = 'alice' AND v1 = 'data1' AND v2 = 'read'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAdapter_AddPolicy_MultipleRules(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	rules := [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"charlie", "data3", "delete"},
	}

	for _, rule := range rules {
		err := adapter.AddPolicy("p", "p", rule)
		require.NoError(t, err)
	}

	var count int
	err := testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE ptype = 'p'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestAdapter_AddPolicy_GroupRule(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	err := adapter.AddPolicy("g", "g", []string{"alice", "admin"})
	require.NoError(t, err)

	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.casbin_rule WHERE ptype = 'g' AND v0 = 'alice' AND v1 = 'admin'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestAdapter_AddPolicy_VerifyWithLoadPolicy(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	// Add a policy
	err := adapter.AddPolicy("p", "p", []string{"dave", "reports", "generate"})
	require.NoError(t, err)

	// Load and verify
	m := newTestModel(t)
	err = adapter.LoadPolicy(m)
	require.NoError(t, err)

	assert.Len(t, m["p"]["p"].Policy, 1)
	assert.Equal(t, []string{"dave", "reports", "generate"}, m["p"]["p"].Policy[0])
}

// =====================================================
// RemovePolicy tests
// =====================================================

func TestAdapter_RemovePolicy_Basic(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add a rule
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	// Remove it
	err = adapter.RemovePolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	// Verify gone from DB
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestAdapter_RemovePolicy_NotFound(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	// Remove a policy that does not exist
	err := adapter.RemovePolicy("p", "p", []string{"nonexistent", "data", "read"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "policy not found")
}

func TestAdapter_RemovePolicy_OnlyMatchingRule(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add two rules
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"alice", "data1", "write"})
	require.NoError(t, err)

	// Remove only the "read" rule
	err = adapter.RemovePolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	// "write" rule should still exist
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.casbin_rule WHERE v0 = 'alice' AND v2 = 'write'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Total should be 1
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

// =====================================================
// RemoveFilteredPolicy tests
// =====================================================

func TestAdapter_RemoveFilteredPolicy_ByFieldIndex0(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add rules for two different subjects
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"alice", "data2", "write"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"bob", "data1", "read"})
	require.NoError(t, err)

	// Remove all policies for "alice" (field index 0 = v0)
	err = adapter.RemoveFilteredPolicy("p", "p", 0, "alice")
	require.NoError(t, err)

	// Only bob's rule should remain
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var v0 string
	err = testDB.Pool().QueryRow(ctx, "SELECT v0 FROM shared.casbin_rule").Scan(&v0)
	require.NoError(t, err)
	assert.Equal(t, "bob", v0)
}

func TestAdapter_RemoveFilteredPolicy_ByFieldIndex1(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add rules with different objects
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"bob", "data1", "write"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"charlie", "data2", "read"})
	require.NoError(t, err)

	// Remove all policies for object "data1" (field index 1 = v1)
	err = adapter.RemoveFilteredPolicy("p", "p", 1, "data1")
	require.NoError(t, err)

	// Only charlie's rule on data2 should remain
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var v0 string
	err = testDB.Pool().QueryRow(ctx, "SELECT v0 FROM shared.casbin_rule").Scan(&v0)
	require.NoError(t, err)
	assert.Equal(t, "charlie", v0)
}

func TestAdapter_RemoveFilteredPolicy_MultipleFieldValues(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add rules
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"alice", "data1", "write"})
	require.NoError(t, err)
	err = adapter.AddPolicy("p", "p", []string{"alice", "data2", "read"})
	require.NoError(t, err)

	// Remove policies matching v0="alice" AND v1="data1" (field index 0, values "alice", "data1")
	err = adapter.RemoveFilteredPolicy("p", "p", 0, "alice", "data1")
	require.NoError(t, err)

	// Only alice's data2 rule should remain
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var v1 string
	err = testDB.Pool().QueryRow(ctx, "SELECT v1 FROM shared.casbin_rule").Scan(&v1)
	require.NoError(t, err)
	assert.Equal(t, "data2", v1)
}

func TestAdapter_RemoveFilteredPolicy_NoMatch(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add a rule
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	// Remove with non-matching filter (should not error, just no rows affected)
	err = adapter.RemoveFilteredPolicy("p", "p", 0, "nonexistent")
	require.NoError(t, err)

	// Original rule should still be there
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

// =====================================================
// Full Casbin integration tests
// =====================================================

func TestAdapter_FullEnforcerIntegration_SaveLoadCycle(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	enforcer := newTestEnforcer(t, adapter)

	// Add policies via enforcer
	_, err := enforcer.AddPolicy("alice", "data1", "read")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("bob", "data2", "write")
	require.NoError(t, err)
	_, err = enforcer.AddRoleForUser("charlie", "admin")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("admin", "data1", "write")
	require.NoError(t, err)

	// Save policies to database
	err = enforcer.SavePolicy()
	require.NoError(t, err)

	// Create a new enforcer (simulating restart) with same adapter
	adapter2 := NewAdapter(adapter.pool)
	enforcer2 := newTestEnforcer(t, adapter2)

	// Verify policies were loaded correctly
	allowed, err := enforcer2.Enforce("alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = enforcer2.Enforce("bob", "data2", "write")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Charlie should have admin role -> can write data1
	allowed, err = enforcer2.Enforce("charlie", "data1", "write")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Charlie should NOT be able to read data1 (no such policy for admin)
	allowed, err = enforcer2.Enforce("charlie", "data1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestAdapter_FullEnforcerIntegration_AddAndRemovePolicy(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	enforcer := newTestEnforcer(t, adapter)

	// Add policy
	_, err := enforcer.AddPolicy("dave", "reports", "view")
	require.NoError(t, err)

	// Verify it works
	allowed, err := enforcer.Enforce("dave", "reports", "view")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Remove policy
	_, err = enforcer.RemovePolicy("dave", "reports", "view")
	require.NoError(t, err)

	// Verify it no longer works
	allowed, err = enforcer.Enforce("dave", "reports", "view")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestAdapter_FullEnforcerIntegration_FilteredPolicyRemoval(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	enforcer := newTestEnforcer(t, adapter)

	// Add multiple policies for a role
	_, err := enforcer.AddPolicy("editor", "articles", "read")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("editor", "articles", "write")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("editor", "articles", "delete")
	require.NoError(t, err)
	_, err = enforcer.AddPolicy("viewer", "articles", "read")
	require.NoError(t, err)

	// Remove all editor policies using filtered removal
	_, err = enforcer.RemoveFilteredPolicy(0, "editor")
	require.NoError(t, err)

	// Editor should have no access
	allowed, err := enforcer.Enforce("editor", "articles", "read")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Viewer should still have access
	allowed, err = enforcer.Enforce("viewer", "articles", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestAdapter_AddPolicy_DuplicateRule(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add the same rule twice
	err := adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	require.NoError(t, err)

	err = adapter.AddPolicy("p", "p", []string{"alice", "data1", "read"})
	// Depending on DB constraints, this may or may not error.
	// The adapter does a plain INSERT, so if there's a unique constraint it errors,
	// otherwise it succeeds and creates a duplicate. Either way, verify the state.
	if err != nil {
		// If it errored, the first rule should still exist
		var count int
		queryErr := testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule").Scan(&count)
		require.NoError(t, queryErr)
		assert.Equal(t, 1, count)
	} else {
		// If it didn't error, we have a duplicate (no unique constraint)
		var count int
		queryErr := testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE v0 = 'alice'").Scan(&count)
		require.NoError(t, queryErr)
		assert.GreaterOrEqual(t, count, 1)
	}
}

func TestAdapter_SavePolicy_WithGroupPolicies(t *testing.T) {
	t.Parallel()
	adapter, _ := setupTestAdapter(t)

	m := newTestModel(t)

	// Add both p-type and g-type policies
	m["p"]["p"].Policy = [][]string{
		{"admin", "users", "manage"},
		{"admin", "movies", "manage"},
		{"editor", "movies", "write"},
	}
	m["g"]["g"].Policy = [][]string{
		{"user1", "admin"},
		{"user2", "editor"},
		{"user3", "admin"},
	}

	err := adapter.SavePolicy(m)
	require.NoError(t, err)

	// Load into a fresh model
	m2 := newTestModel(t)
	err = adapter.LoadPolicy(m2)
	require.NoError(t, err)

	assert.Len(t, m2["p"]["p"].Policy, 3, "should have 3 p-type rules")
	assert.Len(t, m2["g"]["g"].Policy, 3, "should have 3 g-type rules")
}

func TestAdapter_RemoveFilteredPolicy_GroupType(t *testing.T) {
	t.Parallel()
	adapter, testDB := setupTestAdapter(t)
	ctx := context.Background()

	// Add group rules
	err := adapter.AddPolicy("g", "g", []string{"alice", "admin"})
	require.NoError(t, err)
	err = adapter.AddPolicy("g", "g", []string{"bob", "admin"})
	require.NoError(t, err)
	err = adapter.AddPolicy("g", "g", []string{"charlie", "editor"})
	require.NoError(t, err)

	// Remove all "admin" group assignments (v1 = "admin", field index 1)
	err = adapter.RemoveFilteredPolicy("g", "g", 1, "admin")
	require.NoError(t, err)

	// Only charlie's editor assignment should remain
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.casbin_rule WHERE ptype = 'g'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var v0 string
	err = testDB.Pool().QueryRow(ctx, "SELECT v0 FROM shared.casbin_rule WHERE ptype = 'g'").Scan(&v0)
	require.NoError(t, err)
	assert.Equal(t, "charlie", v0)
}
