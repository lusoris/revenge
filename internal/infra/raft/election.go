// Package raft provides leader election for cluster deployments using HashiCorp Raft.
// Used to ensure periodic cleanup jobs run only on the leader node.
package raft

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"log/slog"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
)

// LeaderElection manages Raft-based leader election for a cluster.
// Only the leader should execute periodic cleanup jobs to prevent duplicates.
type LeaderElection struct {
	raft        *raft.Raft
	logger      *slog.Logger
	logStore    *raftboltdb.BoltStore
	stableStore *raftboltdb.BoltStore
	transport   *raft.NetworkTransport
}

// Config holds configuration for Raft leader election.
type Config struct {
	// Enabled controls whether Raft leader election is active
	Enabled bool

	// NodeID is the unique identifier for this node (hostname or UUID)
	NodeID string

	// BindAddr is the address for Raft communication (e.g., "0.0.0.0:7000")
	BindAddr string

	// DataDir is the directory for Raft data storage
	DataDir string

	// Bootstrap should be true only for the first node to initialize the cluster
	Bootstrap bool
}

// NewLeaderElection creates a new Raft-based leader election system.
// Returns nil if Raft is disabled in config.
func NewLeaderElection(cfg Config, logger *slog.Logger) (*LeaderElection, error) {
	if !cfg.Enabled {
		logger.Info("Raft leader election disabled")
		return nil, nil
	}

	logger.Info("Initializing Raft leader election",
		slog.String("node_id", cfg.NodeID),
		slog.String("bind_addr", cfg.BindAddr),
		slog.Bool("bootstrap", cfg.Bootstrap))

	// Setup Raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.NodeID)
	config.Logger = newHCLogAdapter(logger)

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create raft data directory: %w", err)
	}

	// Setup Raft communication
	addr, err := net.ResolveTCPAddr("tcp", cfg.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve bind address: %w", err)
	}

	logWriter := newSlogWriter(logger)

	transport, err := raft.NewTCPTransport(cfg.BindAddr, addr, 3, 10*time.Second, logWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	// Create the snapshot store
	snapshots, err := raft.NewFileSnapshotStore(cfg.DataDir, 2, logWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot store: %w", err)
	}

	// Create the log store and stable store using BoltDB
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-log.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log store: %w", err)
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-stable.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to create stable store: %w", err)
	}

	// Create FSM (simple no-op for leader election only)
	fsm := &simpleFSM{}

	// Instantiate the Raft system
	ra, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create raft instance: %w", err)
	}

	// Bootstrap cluster if this is the first node
	if cfg.Bootstrap {
		logger.Info("Bootstrapping Raft cluster")
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		future := ra.BootstrapCluster(configuration)
		if err := future.Error(); err != nil {
			logger.Warn("Failed to bootstrap cluster (may already be bootstrapped)", slog.Any("error",err))
		}
	}

	return &LeaderElection{
		raft:        ra,
		logger:      logger.With("component", "raft"),
		logStore:    logStore,
		stableStore: stableStore,
		transport:   transport,
	}, nil
}

// IsLeader returns true if this node is the current leader.
func (le *LeaderElection) IsLeader() bool {
	if le == nil || le.raft == nil {
		// If Raft is disabled, consider this node as "leader" (single-node mode)
		return true
	}
	return le.raft.State() == raft.Leader
}

// LeaderAddr returns the address of the current leader.
func (le *LeaderElection) LeaderAddr() string {
	if le == nil || le.raft == nil {
		return "single-node"
	}
	addr, id := le.raft.LeaderWithID()
	if addr == "" {
		return "no-leader"
	}
	return fmt.Sprintf("%s (%s)", addr, id)
}

// State returns the current Raft state as a string.
func (le *LeaderElection) State() string {
	if le == nil || le.raft == nil {
		return "disabled"
	}
	switch le.raft.State() {
	case raft.Leader:
		return "leader"
	case raft.Candidate:
		return "candidate"
	case raft.Follower:
		return "follower"
	case raft.Shutdown:
		return "shutdown"
	default:
		return "unknown"
	}
}

// Close shuts down the Raft instance gracefully and releases all resources.
func (le *LeaderElection) Close() error {
	if le == nil || le.raft == nil {
		return nil
	}
	le.logger.Info("Shutting down Raft")
	if err := le.raft.Shutdown().Error(); err != nil {
		return err
	}
	if le.transport != nil {
		if err := le.transport.Close(); err != nil {
			le.logger.Warn("Failed to close Raft transport", slog.Any("error", err))
		}
	}
	if le.logStore != nil {
		if err := le.logStore.Close(); err != nil {
			le.logger.Warn("Failed to close Raft log store", slog.Any("error", err))
		}
	}
	if le.stableStore != nil {
		if err := le.stableStore.Close(); err != nil {
			le.logger.Warn("Failed to close Raft stable store", slog.Any("error", err))
		}
	}
	return nil
}

// AddVoter adds a new node to the Raft cluster as a voter.
// Must be called on the leader node. The joining node should already be
// running and listening on the specified address.
func (le *LeaderElection) AddVoter(nodeID, addr string) error {
	if le == nil || le.raft == nil {
		return fmt.Errorf("raft is not initialized")
	}
	if le.raft.State() != raft.Leader {
		return fmt.Errorf("not the leader, cannot add voter")
	}

	le.logger.Info("Adding voter to Raft cluster",
		slog.String("node_id", nodeID),
		slog.String("address", addr),
	)

	future := le.raft.AddVoter(
		raft.ServerID(nodeID),
		raft.ServerAddress(addr),
		0, // prevIndex=0 lets Raft handle log index
		30*time.Second,
	)
	if err := future.Error(); err != nil {
		return fmt.Errorf("failed to add voter %s: %w", nodeID, err)
	}

	le.logger.Info("Voter added to Raft cluster",
		slog.String("node_id", nodeID),
		slog.String("address", addr),
	)
	return nil
}

// RemoveServer removes a node from the Raft cluster.
// Must be called on the leader node.
func (le *LeaderElection) RemoveServer(nodeID string) error {
	if le == nil || le.raft == nil {
		return fmt.Errorf("raft is not initialized")
	}
	if le.raft.State() != raft.Leader {
		return fmt.Errorf("not the leader, cannot remove server")
	}

	le.logger.Info("Removing server from Raft cluster",
		slog.String("node_id", nodeID),
	)

	future := le.raft.RemoveServer(
		raft.ServerID(nodeID),
		0, // prevIndex=0 lets Raft handle log index
		30*time.Second,
	)
	if err := future.Error(); err != nil {
		return fmt.Errorf("failed to remove server %s: %w", nodeID, err)
	}

	le.logger.Info("Server removed from Raft cluster",
		slog.String("node_id", nodeID),
	)
	return nil
}

// ClusterMember represents a node in the Raft cluster.
type ClusterMember struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Voter   bool   `json:"voter"`
}

// GetClusterMembers returns the current Raft cluster configuration.
func (le *LeaderElection) GetClusterMembers() ([]ClusterMember, error) {
	if le == nil || le.raft == nil {
		return []ClusterMember{{ID: "single-node", Address: "local", Voter: true}}, nil
	}

	future := le.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, fmt.Errorf("failed to get cluster configuration: %w", err)
	}

	servers := future.Configuration().Servers
	members := make([]ClusterMember, 0, len(servers))
	for _, server := range servers {
		members = append(members, ClusterMember{
			ID:      string(server.ID),
			Address: string(server.Address),
			Voter:   server.Suffrage == raft.Voter,
		})
	}
	return members, nil
}

// simpleFSM is a minimal FSM implementation for leader election.
// We don't need to maintain any state - just need leader election.
type simpleFSM struct{}

func (f *simpleFSM) Apply(*raft.Log) interface{} {
	// No state to apply
	return nil
}

func (f *simpleFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &simpleSnapshot{}, nil
}

func (f *simpleFSM) Restore(rc io.ReadCloser) error {
	// No state to restore
	return rc.Close()
}

// simpleSnapshot is a minimal snapshot implementation.
type simpleSnapshot struct{}

func (s *simpleSnapshot) Persist(sink raft.SnapshotSink) error {
	// No state to persist
	return sink.Close()
}

func (s *simpleSnapshot) Release() {
	// Nothing to release
}

// hcLogAdapter adapts slog.Logger to hashicorp/go-hclog interface.
type hcLogAdapter struct {
	logger *slog.Logger
	level  hclog.Level
}

func newHCLogAdapter(logger *slog.Logger) hclog.Logger {
	return &hcLogAdapter{
		logger: logger.With("component", "raft"),
		level:  hclog.Info,
	}
}

func (h *hcLogAdapter) Log(level hclog.Level, msg string, args ...interface{}) {
	attrs := hclogArgsToAttrs(args)
	switch level {
	case hclog.Trace, hclog.Debug:
		h.logger.Debug(msg, attrs...)
	case hclog.Info:
		h.logger.Info(msg, attrs...)
	case hclog.Warn:
		h.logger.Warn(msg, attrs...)
	case hclog.Error:
		h.logger.Error(msg, attrs...)
	}
}

func (h *hcLogAdapter) Trace(msg string, args ...interface{}) {
	h.logger.Debug(msg, hclogArgsToAttrs(args)...)
}
func (h *hcLogAdapter) Debug(msg string, args ...interface{}) {
	h.logger.Debug(msg, hclogArgsToAttrs(args)...)
}
func (h *hcLogAdapter) Info(msg string, args ...interface{}) {
	h.logger.Info(msg, hclogArgsToAttrs(args)...)
}
func (h *hcLogAdapter) Warn(msg string, args ...interface{}) {
	h.logger.Warn(msg, hclogArgsToAttrs(args)...)
}
func (h *hcLogAdapter) Error(msg string, args ...interface{}) {
	h.logger.Error(msg, hclogArgsToAttrs(args)...)
}

func (h *hcLogAdapter) IsTrace() bool { return h.level <= hclog.Trace }
func (h *hcLogAdapter) IsDebug() bool { return h.level <= hclog.Debug }
func (h *hcLogAdapter) IsInfo() bool  { return h.level <= hclog.Info }
func (h *hcLogAdapter) IsWarn() bool  { return h.level <= hclog.Warn }
func (h *hcLogAdapter) IsError() bool { return h.level <= hclog.Error }

func (h *hcLogAdapter) ImpliedArgs() []interface{}        { return nil }
func (h *hcLogAdapter) With(args ...interface{}) hclog.Logger { return h }
func (h *hcLogAdapter) Name() string                      { return "raft" }
func (h *hcLogAdapter) Named(name string) hclog.Logger    { return h }
func (h *hcLogAdapter) ResetNamed(name string) hclog.Logger { return h }
func (h *hcLogAdapter) SetLevel(level hclog.Level)        { h.level = level }
func (h *hcLogAdapter) GetLevel() hclog.Level             { return h.level }
func (h *hcLogAdapter) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}
func (h *hcLogAdapter) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return newSlogWriter(h.logger)
}

// hclogArgsToAttrs converts hashicorp/go-hclog key-value pairs to slog attributes.
// go-hclog passes args as alternating key, value pairs: "key1", val1, "key2", val2, ...
func hclogArgsToAttrs(args []interface{}) []any {
	if len(args) == 0 {
		return nil
	}
	attrs := make([]any, 0, len(args))
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", args[i])
		}
		attrs = append(attrs, slog.Any(key, args[i+1]))
	}
	// Handle odd number of args (trailing key without value)
	if len(args)%2 != 0 {
		attrs = append(attrs, slog.Any("EXTRA_VALUE_AT_END", args[len(args)-1]))
	}
	return attrs
}

// slogWriter adapts slog.Logger to io.Writer for Raft transport and snapshot logging.
type slogWriter struct {
	logger *slog.Logger
}

func newSlogWriter(logger *slog.Logger) io.Writer {
	return &slogWriter{logger: logger.With("component", "raft-transport")}
}

func (w *slogWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	// Trim trailing newline
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	w.logger.Debug(msg)
	return len(p), nil
}
