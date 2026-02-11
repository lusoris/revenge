package raft

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testSlogLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestLeaderElection_IsLeader_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	assert.True(t, le.IsLeader(), "nil LeaderElection should be considered leader (single-node)")
}

func TestLeaderElection_IsLeader_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	assert.True(t, le.IsLeader(), "LeaderElection with nil raft should be considered leader")
}

func TestLeaderElection_LeaderAddr_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	assert.Equal(t, "single-node", le.LeaderAddr())
}

func TestLeaderElection_LeaderAddr_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	assert.Equal(t, "single-node", le.LeaderAddr())
}

func TestLeaderElection_State_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	assert.Equal(t, "disabled", le.State())
}

func TestLeaderElection_State_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	assert.Equal(t, "disabled", le.State())
}

func TestLeaderElection_Close_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	assert.NoError(t, le.Close())
}

func TestLeaderElection_Close_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	assert.NoError(t, le.Close())
}

func TestLeaderElection_AddVoter_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	err := le.AddVoter("node-2", "127.0.0.1:7001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "raft is not initialized")
}

func TestLeaderElection_AddVoter_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	err := le.AddVoter("node-2", "127.0.0.1:7001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "raft is not initialized")
}

func TestLeaderElection_RemoveServer_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	err := le.RemoveServer("node-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "raft is not initialized")
}

func TestLeaderElection_RemoveServer_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	err := le.RemoveServer("node-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "raft is not initialized")
}

func TestLeaderElection_GetClusterMembers_Nil(t *testing.T) {
	t.Parallel()
	var le *LeaderElection
	members, err := le.GetClusterMembers()
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, "single-node", members[0].ID)
	assert.Equal(t, "local", members[0].Address)
	assert.True(t, members[0].Voter)
}

func TestLeaderElection_GetClusterMembers_NilRaft(t *testing.T) {
	t.Parallel()
	le := &LeaderElection{}
	members, err := le.GetClusterMembers()
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, "single-node", members[0].ID)
}

func TestNewLeaderElection_Disabled(t *testing.T) {
	t.Parallel()
	cfg := Config{
		Enabled: false,
	}
	le, err := NewLeaderElection(cfg, testSlogLogger())
	require.NoError(t, err)
	assert.Nil(t, le)
}

func TestClusterMember_Struct(t *testing.T) {
	t.Parallel()
	m := ClusterMember{
		ID:      "node-1",
		Address: "10.0.0.1:7000",
		Voter:   true,
	}
	assert.Equal(t, "node-1", m.ID)
	assert.Equal(t, "10.0.0.1:7000", m.Address)
	assert.True(t, m.Voter)
}
