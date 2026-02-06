package raft

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides Raft leader election for fx dependency injection.
var Module = fx.Module("raft",
	fx.Provide(provideLeaderElection),
	fx.Invoke(registerLifecycle),
)

// provideLeaderElection creates a Raft leader election instance.
// Returns nil if Raft is disabled (single-node mode).
func provideLeaderElection(cfg *config.Config, logger *zap.Logger) (*LeaderElection, error) {
	// Generate node ID if not provided
	nodeID := cfg.Raft.NodeID
	if nodeID == "" {
		// Use hostname as default node ID
		hostname, err := os.Hostname()
		if err != nil {
			// Fallback to UUID
			nodeID = uuid.Must(uuid.NewV7()).String()
			logger.Warn("Failed to get hostname, using UUID as node ID",
				zap.String("node_id", nodeID),
				zap.Error(err))
		} else {
			nodeID = hostname
		}
	}

	raftConfig := Config{
		Enabled:   cfg.Raft.Enabled,
		NodeID:    nodeID,
		BindAddr:  cfg.Raft.BindAddr,
		DataDir:   cfg.Raft.DataDir,
		Bootstrap: cfg.Raft.Bootstrap,
	}

	return NewLeaderElection(raftConfig, logger)
}

// registerLifecycle registers Raft lifecycle hooks with fx.
func registerLifecycle(lc fx.Lifecycle, le *LeaderElection, logger *zap.Logger) {
	if le == nil {
		// Raft disabled
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Raft leader election started",
				zap.String("state", le.State()),
				zap.String("leader", le.LeaderAddr()))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Raft leader election")
			return le.Close()
		},
	})
}
