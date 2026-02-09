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
)

// providePeriodicJobs builds the list of periodic jobs for the River client.
func providePeriodicJobs(cfg *config.Config) []*river.PeriodicJob {
	retentionDays := cfg.Activity.RetentionDays
	if retentionDays <= 0 {
		retentionDays = 90
	}

	return []*river.PeriodicJob{
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
	}
}

// registerActivityCleanupWorker registers the activity cleanup worker with River.
func registerActivityCleanupWorker(workers *river.Workers, worker *activity.ActivityCleanupWorker) {
	river.AddWorker(workers, worker)
}

// registerLibraryCleanupWorker registers the library scan cleanup worker with River.
func registerLibraryCleanupWorker(workers *river.Workers, worker *library.LibraryScanCleanupWorker) {
	river.AddWorker(workers, worker)
}
