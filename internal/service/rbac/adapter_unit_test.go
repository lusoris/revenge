package rbac

import (
	"testing"

	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for adapter.go functions that do not require a database connection.

func TestLoadPolicyLine_AllFieldsPopulated(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	rule := &CasbinRule{
		PType: "p",
		V0:    "admin",
		V1:    "users",
		V2:    "manage",
		V3:    "domain1",
		V4:    "extra4",
		V5:    "extra5",
	}

	loadPolicyLine(rule, m)

	policies := m["p"]["p"].Policy
	require.Len(t, policies, 1)
	assert.Contains(t, policies[0], "admin")
	assert.Contains(t, policies[0], "users")
	assert.Contains(t, policies[0], "manage")
}

func TestLoadPolicyLine_OnlyV0(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	rule := &CasbinRule{
		PType: "p",
		V0:    "alice",
		V1:    "",
		V2:    "",
	}

	loadPolicyLine(rule, m)

	policies := m["p"]["p"].Policy
	require.Len(t, policies, 1)
	assert.Equal(t, []string{"alice"}, policies[0])
}

func TestLoadPolicyLine_GTypeWithV0V1(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	rule := &CasbinRule{
		PType: "g",
		V0:    "bob",
		V1:    "moderator",
	}

	loadPolicyLine(rule, m)

	groupPolicies := m["g"]["g"].Policy
	require.Len(t, groupPolicies, 1)
	assert.Equal(t, []string{"bob", "moderator"}, groupPolicies[0])
}

func TestLoadPolicyLine_UnknownSectionIgnored(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	// "z" is not a valid section, should be silently ignored
	rule := &CasbinRule{
		PType: "z",
		V0:    "alice",
		V1:    "data",
		V2:    "read",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		loadPolicyLine(rule, m)
	})
}

func TestLoadPolicyLine_UnknownKeyInKnownSection(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	// "p2" key doesn't exist in the model under section "p"
	rule := &CasbinRule{
		PType: "p2",
		V0:    "alice",
		V1:    "data",
		V2:    "read",
	}

	// Should not panic - section "p" exists but key "p2" does not
	assert.NotPanics(t, func() {
		loadPolicyLine(rule, m)
	})
}

func TestLoadPolicyLine_PartialV3V4V5(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	// V3 populated, V4 empty, V5 populated
	rule := &CasbinRule{
		PType: "p",
		V0:    "alice",
		V1:    "data1",
		V2:    "read",
		V3:    "extra3",
		V4:    "",
		V5:    "extra5",
	}

	loadPolicyLine(rule, m)

	policies := m["p"]["p"].Policy
	require.Len(t, policies, 1)
	// V3="extra3" should be appended, V4="" skipped, V5="extra5" should NOT be appended
	// because loadPolicyLine stops at the first empty field
	assert.Contains(t, policies[0], "alice")
	assert.Contains(t, policies[0], "extra3")
}

func TestNewAdapter_Properties(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter(nil)

	assert.NotNil(t, adapter)
	assert.Nil(t, adapter.pool)
	assert.Equal(t, "shared.casbin_rule", adapter.tableName)
}

func TestCasbinRule_EmptyStruct(t *testing.T) {
	t.Parallel()

	rule := CasbinRule{}

	assert.Empty(t, rule.PType)
	assert.Empty(t, rule.V0)
	assert.Empty(t, rule.V1)
	assert.Empty(t, rule.V2)
	assert.Empty(t, rule.V3)
	assert.Empty(t, rule.V4)
	assert.Empty(t, rule.V5)
}

func TestLoadPolicyLine_BatchLoadAndVerify(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	rules := []*CasbinRule{
		{PType: "p", V0: "admin", V1: "users", V2: "manage"},
		{PType: "p", V0: "admin", V1: "movies", V2: "delete"},
		{PType: "p", V0: "editor", V1: "movies", V2: "write"},
		{PType: "g", V0: "alice", V1: "admin"},
		{PType: "g", V0: "bob", V1: "editor"},
	}

	for _, rule := range rules {
		loadPolicyLine(rule, m)
	}

	assert.Len(t, m["p"]["p"].Policy, 3, "should have 3 policy rules")
	assert.Len(t, m["g"]["g"].Policy, 2, "should have 2 group rules")
}
