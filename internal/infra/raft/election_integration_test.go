//go:build integration

package raft

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getFreePort returns a random available TCP port.
func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// createNode boots a single Raft node with a temp directory and random port.
func createNode(t *testing.T, nodeID string, bootstrap bool) (*LeaderElection, Config) {
	t.Helper()

	port := getFreePort(t)
	dataDir := t.TempDir()

	cfg := Config{
		Enabled:   true,
		NodeID:    nodeID,
		BindAddr:  fmt.Sprintf("127.0.0.1:%d", port),
		DataDir:   dataDir,
		Bootstrap: bootstrap,
	}

	le, err := NewLeaderElection(cfg, testLogger())
	require.NoError(t, err)
	require.NotNil(t, le)

	t.Cleanup(func() {
		if le != nil {
			le.Close() //nolint:errcheck
		}
	})

	return le, cfg
}

// waitForLeader polls until the node reports itself as leader or times out.
func waitForLeader(t *testing.T, le *LeaderElection, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if le.IsLeader() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("node did not become leader within %v (state: %s)", timeout, le.State())
}

// waitForState polls until the node reaches the expected state or times out.
func waitForState(t *testing.T, le *LeaderElection, expected string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if le.State() == expected {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("node did not reach state %q within %v (state: %s)", expected, timeout, le.State())
}

// findLeader returns the first node that reports as leader.
func findLeader(nodes []*LeaderElection) *LeaderElection {
	for _, n := range nodes {
		if n.IsLeader() {
			return n
		}
	}
	return nil
}

// ============================================================================
// Test 1: Single-node bootstrap
// Boot one node, verify it elects itself leader
// ============================================================================

func TestRaft_SingleNodeBootstrap(t *testing.T) {
	le, _ := createNode(t, "node-1", true)

	// Should become leader quickly (single-node cluster)
	waitForLeader(t, le, 10*time.Second)

	t.Run("is_leader", func(t *testing.T) {
		assert.True(t, le.IsLeader())
	})

	t.Run("state_is_leader", func(t *testing.T) {
		assert.Equal(t, "leader", le.State())
	})

	t.Run("leader_addr_not_empty", func(t *testing.T) {
		addr := le.LeaderAddr()
		assert.NotEmpty(t, addr)
		assert.NotEqual(t, "no-leader", addr)
		assert.NotEqual(t, "single-node", addr)
	})

	t.Run("cluster_members", func(t *testing.T) {
		members, err := le.GetClusterMembers()
		require.NoError(t, err)
		require.Len(t, members, 1)
		assert.Equal(t, "node-1", members[0].ID)
		assert.True(t, members[0].Voter)
	})

	t.Run("graceful_close", func(t *testing.T) {
		err := le.Close()
		assert.NoError(t, err)
		assert.Equal(t, "shutdown", le.State())
	})
}

// ============================================================================
// Test 2: Multi-node cluster (3 nodes)
// Bootstrap one node, join two more, verify cluster membership
// ============================================================================

func TestRaft_MultiNodeCluster(t *testing.T) {
	// Boot the bootstrap node (node-1)
	node1, cfg1 := createNode(t, "node-1", true)
	waitForLeader(t, node1, 10*time.Second)

	// Boot node-2 (non-bootstrap)
	node2, cfg2 := createNode(t, "node-2", false)

	// Boot node-3 (non-bootstrap)
	node3, cfg3 := createNode(t, "node-3", false)

	// Leader (node-1) adds node-2 and node-3 as voters
	t.Run("add_voter_node2", func(t *testing.T) {
		err := node1.AddVoter("node-2", cfg2.BindAddr)
		require.NoError(t, err)
	})

	t.Run("add_voter_node3", func(t *testing.T) {
		err := node1.AddVoter("node-3", cfg3.BindAddr)
		require.NoError(t, err)
	})

	// Wait for cluster to stabilize
	time.Sleep(2 * time.Second)

	// Verify cluster membership from leader
	t.Run("cluster_has_3_members", func(t *testing.T) {
		members, err := node1.GetClusterMembers()
		require.NoError(t, err)
		assert.Len(t, members, 3, "cluster should have 3 members")

		ids := make(map[string]bool)
		for _, m := range members {
			ids[m.ID] = true
			assert.True(t, m.Voter, "all members should be voters")
		}
		assert.True(t, ids["node-1"])
		assert.True(t, ids["node-2"])
		assert.True(t, ids["node-3"])
	})

	// Exactly one leader
	t.Run("exactly_one_leader", func(t *testing.T) {
		leaders := 0
		for _, n := range []*LeaderElection{node1, node2, node3} {
			if n.IsLeader() {
				leaders++
			}
		}
		assert.Equal(t, 1, leaders, "cluster should have exactly one leader")
	})

	// Followers know the leader address
	t.Run("followers_know_leader", func(t *testing.T) {
		for _, n := range []*LeaderElection{node2, node3} {
			addr := n.LeaderAddr()
			assert.NotEqual(t, "no-leader", addr, "followers should know leader")
		}
	})

	// Remove node-3
	t.Run("remove_server", func(t *testing.T) {
		leader := findLeader([]*LeaderElection{node1, node2, node3})
		require.NotNil(t, leader, "must have a leader to remove a server")

		err := leader.RemoveServer("node-3")
		require.NoError(t, err)

		time.Sleep(1 * time.Second)

		members, err := leader.GetClusterMembers()
		require.NoError(t, err)
		assert.Len(t, members, 2, "cluster should have 2 members after removal")
	})

	_ = cfg1
}

// ============================================================================
// Test 3: Leader failover
// 3-node cluster → shut down leader → verify new leader elected
// ============================================================================

func TestRaft_LeaderFailover(t *testing.T) {
	// Boot 3-node cluster
	node1, _ := createNode(t, "failover-1", true)
	waitForLeader(t, node1, 10*time.Second)

	node2, cfg2 := createNode(t, "failover-2", false)
	node3, cfg3 := createNode(t, "failover-3", false)

	require.NoError(t, node1.AddVoter("failover-2", cfg2.BindAddr))
	require.NoError(t, node1.AddVoter("failover-3", cfg3.BindAddr))

	// Wait for cluster to stabilize
	time.Sleep(3 * time.Second)

	// Identify the leader
	nodes := []*LeaderElection{node1, node2, node3}
	var leader *LeaderElection
	var followers []*LeaderElection
	for _, n := range nodes {
		if n.IsLeader() {
			leader = n
		} else {
			followers = append(followers, n)
		}
	}
	require.NotNil(t, leader, "must have an initial leader")
	require.Len(t, followers, 2, "must have 2 followers")

	oldLeaderState := leader.State()
	assert.Equal(t, "leader", oldLeaderState)

	// Shut down the leader
	t.Run("shutdown_leader", func(t *testing.T) {
		err := leader.Close()
		require.NoError(t, err)
	})

	// Wait for a new leader to emerge from the followers
	t.Run("new_leader_elected", func(t *testing.T) {
		deadline := time.Now().Add(15 * time.Second)
		var newLeader *LeaderElection
		for time.Now().Before(deadline) {
			for _, f := range followers {
				if f.IsLeader() {
					newLeader = f
					break
				}
			}
			if newLeader != nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		require.NotNil(t, newLeader, "a follower should become the new leader")
		assert.Equal(t, "leader", newLeader.State())
	})

	// Cluster still has members (from the perspective of the new leader)
	t.Run("cluster_still_functional", func(t *testing.T) {
		newLeader := findLeader(followers)
		require.NotNil(t, newLeader)

		members, err := newLeader.GetClusterMembers()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(members), 2, "cluster should still have members")
	})
}

// ============================================================================
// Test 4: Non-leader cannot add/remove voters
// ============================================================================

func TestRaft_NonLeaderCannotMutateCluster(t *testing.T) {
	node1, _ := createNode(t, "noleader-1", true)
	waitForLeader(t, node1, 10*time.Second)

	node2, cfg2 := createNode(t, "noleader-2", false)
	require.NoError(t, node1.AddVoter("noleader-2", cfg2.BindAddr))

	time.Sleep(2 * time.Second)

	// node2 should be a follower
	t.Run("follower_cannot_add_voter", func(t *testing.T) {
		if node2.IsLeader() {
			t.Skip("node2 became leader unexpectedly")
		}
		err := node2.AddVoter("noleader-3", "127.0.0.1:9999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not the leader")
	})

	t.Run("follower_cannot_remove_server", func(t *testing.T) {
		if node2.IsLeader() {
			t.Skip("node2 became leader unexpectedly")
		}
		err := node2.RemoveServer("noleader-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not the leader")
	})
}

// ============================================================================
// Test 5: Disabled Raft returns nil and behaves as single-node
// ============================================================================

func TestRaft_DisabledReturnsNil(t *testing.T) {
	cfg := Config{Enabled: false}
	le, err := NewLeaderElection(cfg, testLogger())
	require.NoError(t, err)
	assert.Nil(t, le)

	// Nil LeaderElection should act as single-node leader
	assert.True(t, le.IsLeader())
	assert.Equal(t, "single-node", le.LeaderAddr())
	assert.Equal(t, "disabled", le.State())
	assert.NoError(t, le.Close())
}

// ============================================================================
// Test 6: Data directory creation
// ============================================================================

func TestRaft_CreatesDataDirectory(t *testing.T) {
	port := getFreePort(t)
	dataDir := t.TempDir() + "/nested/raft/data"

	cfg := Config{
		Enabled:   true,
		NodeID:    "dir-test-node",
		BindAddr:  fmt.Sprintf("127.0.0.1:%d", port),
		DataDir:   dataDir,
		Bootstrap: true,
	}

	le, err := NewLeaderElection(cfg, testLogger())
	require.NoError(t, err)
	require.NotNil(t, le)
	defer le.Close() //nolint:errcheck

	// Verify directory was created
	info, err := os.Stat(dataDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Verify BoltDB file exists
	_, err = os.Stat(dataDir + "/raft.db")
	assert.NoError(t, err, "raft.db should exist in data dir")
}

// ============================================================================
// Test 7: Invalid bind address
// ============================================================================

func TestRaft_InvalidBindAddress(t *testing.T) {
	cfg := Config{
		Enabled:  true,
		NodeID:   "bad-addr",
		BindAddr: "not-a-valid-address",
		DataDir:  t.TempDir(),
	}

	le, err := NewLeaderElection(cfg, testLogger())
	assert.Error(t, err)
	assert.Nil(t, le)
}

// ============================================================================
// Test 8: FSM and Snapshot are no-ops
// ============================================================================

func TestRaft_SimpleFSM(t *testing.T) {
	fsm := &simpleFSM{}

	t.Run("apply_returns_nil", func(t *testing.T) {
		result := fsm.Apply(nil)
		assert.Nil(t, result)
	})

	t.Run("snapshot_succeeds", func(t *testing.T) {
		snap, err := fsm.Snapshot()
		require.NoError(t, err)
		assert.NotNil(t, snap)
	})

	t.Run("restore_closes_reader", func(t *testing.T) {
		r := io.NopCloser(nil)
		err := fsm.Restore(r)
		assert.NoError(t, err)
	})
}

// ============================================================================
// Test 9: HCLog adapter
// ============================================================================

func TestRaft_HCLogAdapter(t *testing.T) {
	adapter := newHCLogAdapter(testLogger())

	t.Run("level_defaults_info", func(t *testing.T) {
		assert.True(t, adapter.IsInfo())
		assert.False(t, adapter.IsDebug())
	})

	t.Run("name", func(t *testing.T) {
		assert.Equal(t, "raft", adapter.Name())
	})

	t.Run("log_methods_dont_panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			adapter.Trace("trace msg", "key", "val")
			adapter.Debug("debug msg", "key", "val")
			adapter.Info("info msg", "key", "val")
			adapter.Warn("warn msg", "key", "val")
			adapter.Error("error msg", "key", "val")
		})
	})

	t.Run("with_returns_adapter", func(t *testing.T) {
		child := adapter.With("extra", "arg")
		assert.NotNil(t, child)
	})

	t.Run("named_returns_adapter", func(t *testing.T) {
		named := adapter.Named("sub")
		assert.NotNil(t, named)
	})

	t.Run("standard_logger", func(t *testing.T) {
		stdLogger := adapter.StandardLogger(nil)
		assert.NotNil(t, stdLogger)
	})

	t.Run("standard_writer", func(t *testing.T) {
		w := adapter.StandardWriter(nil)
		assert.NotNil(t, w)
		n, err := w.Write([]byte("test\n"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
	})
}

// ============================================================================
// Test 10: hclogArgsToAttrs edge cases
// ============================================================================

func TestRaft_HCLogArgsToAttrs(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		attrs := hclogArgsToAttrs(nil)
		assert.Nil(t, attrs)
	})

	t.Run("key_value_pairs", func(t *testing.T) {
		attrs := hclogArgsToAttrs([]interface{}{"key1", "val1", "key2", 42})
		assert.Len(t, attrs, 2)
	})

	t.Run("odd_trailing_value", func(t *testing.T) {
		attrs := hclogArgsToAttrs([]interface{}{"key1", "val1", "orphan"})
		assert.Len(t, attrs, 2) // key1->val1 + EXTRA_VALUE_AT_END->orphan
	})

	t.Run("non_string_key", func(t *testing.T) {
		attrs := hclogArgsToAttrs([]interface{}{123, "val"})
		assert.Len(t, attrs, 1)
	})
}
