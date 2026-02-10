package analytics

import (
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

// Module provides analytics service dependencies.
var Module = fx.Module("analytics",
	fx.Provide(NewStatsAggregationWorker),
	fx.Invoke(registerStatsAggregationWorker),
)

// registerStatsAggregationWorker registers the worker with River.
func registerStatsAggregationWorker(workers *river.Workers, worker *StatsAggregationWorker) {
	river.AddWorker(workers, worker)
}
