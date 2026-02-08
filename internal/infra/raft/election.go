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

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"log/slog"
)

// LeaderElection manages Raft-based leader election for a cluster.
// Only the leader should execute periodic cleanup jobs to prevent duplicates.
type LeaderElection struct {
	raft   *raft.Raft
	logger *slog.Logger
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

	transport, err := raft.NewTCPTransport(cfg.BindAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	// Create the snapshot store
	snapshots, err := raft.NewFileSnapshotStore(cfg.DataDir, 2, os.Stderr)
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
		raft:   ra,
		logger: logger.With("component", "raft"),
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

// Close shuts down the Raft instance gracefully.
func (le *LeaderElection) Close() error {
	if le == nil || le.raft == nil {
		return nil
	}
	le.logger.Info("Shutting down Raft")
	return le.raft.Shutdown().Error()
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
	switch level {
	case hclog.Trace, hclog.Debug:
		h.logger.Debug(msg)
	case hclog.Info:
		h.logger.Info(msg)
	case hclog.Warn:
		h.logger.Warn(msg)
	case hclog.Error:
		h.logger.Error(msg)
	}
}

func (h *hcLogAdapter) Trace(msg string, args ...interface{}) { h.logger.Debug(msg) }
func (h *hcLogAdapter) Debug(msg string, args ...interface{}) { h.logger.Debug(msg) }
func (h *hcLogAdapter) Info(msg string, args ...interface{})  { h.logger.Info(msg) }
func (h *hcLogAdapter) Warn(msg string, args ...interface{})  { h.logger.Warn(msg) }
func (h *hcLogAdapter) Error(msg string, args ...interface{}) { h.logger.Error(msg) }

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
	return os.Stderr
}
