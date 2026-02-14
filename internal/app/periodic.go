// Package app provides the main application module that wires all dependencies together.
package app

import (
	"time"

	"github.com/riverqueue/river"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/jobs"
	playbackjobs "github.com/lusoris/revenge/internal/playback/jobs"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/analytics"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/session"
)

// providePeriodicJobs builds the list of periodic jobs for the River client.
func providePeriodicJobs(cfg *config.Config) []*river.PeriodicJob {
	retentionDays := cfg.Activity.RetentionDays
	if retentionDays <= 0 {
		retentionDays = 90
	}

	periodicJobs := []*river.PeriodicJob{
		// Auth cleanup: remove expired tokens, password resets, etc. (daily)
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupArgs{
					TargetType: jobs.CleanupTargetAll,
					OlderThan:  24 * time.Hour,
				}, nil
			},
			&river.PeriodicJobOpts{ID: "auth_cleanup_daily", RunOnStart: true},
		),

		// Activity cleanup: remove old activity logs (daily)
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return activity.ActivityCleanupArgs{
					RetentionDays: retentionDays,
				}, nil
			},
			&river.PeriodicJobOpts{ID: "activity_cleanup_daily"},
		),

		// Library scan cleanup: remove old scan history (daily)
		river.NewPeriodicJob(
			river.PeriodicInterval(24*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return library.LibraryScanCleanupArgs{
					RetentionDays: 30,
				}, nil
			},
			&river.PeriodicJobOpts{ID: "library_scan_cleanup_daily"},
		),

		// Playback health check: log session health (every 5 minutes)
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return playbackjobs.CleanupArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "playback_health_check"},
		),

		// Stats aggregation: compute server-wide metrics (hourly)
		river.NewPeriodicJob(
			river.PeriodicInterval(1*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return analytics.StatsAggregationArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "stats_aggregation_hourly", RunOnStart: true},
		),

		// Session maintenance: clean up expired/revoked sessions + reconcile gauge (hourly)
		river.NewPeriodicJob(
			river.PeriodicInterval(1*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return session.MaintenanceArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "session_maintenance_hourly", RunOnStart: true},
		),
	}

	// Periodic library scan: auto-scan all enabled libraries on a schedule.
	// Only enabled when ScanInterval > 0 (default "0s" = disabled).
	if cfg.Movie.Library.ScanInterval > 0 {
		periodicJobs = append(periodicJobs, river.NewPeriodicJob(
			river.PeriodicInterval(cfg.Movie.Library.ScanInterval),
			func() (river.JobArgs, *river.InsertOpts) {
				return library.PeriodicLibraryScanArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "periodic_library_scan"},
		))
	}

	return periodicJobs
}

// registerActivityCleanupWorker registers the activity cleanup worker with River.
func registerActivityCleanupWorker(workers *river.Workers, worker *activity.ActivityCleanupWorker) {
	river.AddWorker(workers, worker)
}

// registerLibraryCleanupWorker registers the library scan cleanup worker with River.
func registerLibraryCleanupWorker(workers *river.Workers, worker *library.LibraryScanCleanupWorker) {
	river.AddWorker(workers, worker)
}

// registerPeriodicLibraryScanWorker registers the periodic library scan worker with River.
func registerPeriodicLibraryScanWorker(workers *river.Workers, worker *library.PeriodicLibraryScanWorker) {
	river.AddWorker(workers, worker)
}
