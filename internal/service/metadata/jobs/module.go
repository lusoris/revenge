package jobs

import (
	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

// Module provides the metadata job workers and queue via fx.
var Module = fx.Module("metadatajobs",
	fx.Provide(
		NewQueue,
		NewRefreshTVShowWorker,
		NewRefreshSeasonWorker,
		NewRefreshEpisodeWorker,
		NewRefreshPersonWorker,
		NewEnrichContentWorker,
		NewDownloadImageWorker,
	),
	fx.Invoke(RegisterWorkers, WireJobQueue),
)

// RegisterWorkers registers all metadata job workers with the River workers registry.
// Note: RefreshMovieArgs worker is already registered by moviejobs.Module —
// we only register workers for the 6 previously orphaned job kinds.
func RegisterWorkers(
	workers *river.Workers,
	tvshowWorker *RefreshTVShowWorker,
	seasonWorker *RefreshSeasonWorker,
	episodeWorker *RefreshEpisodeWorker,
	personWorker *RefreshPersonWorker,
	enrichWorker *EnrichContentWorker,
	imageWorker *DownloadImageWorker,
) error {
	river.AddWorker(workers, tvshowWorker)
	river.AddWorker(workers, seasonWorker)
	river.AddWorker(workers, episodeWorker)
	river.AddWorker(workers, personWorker)
	river.AddWorker(workers, enrichWorker)
	river.AddWorker(workers, imageWorker)
	return nil
}

// WireJobQueue connects the metadata Queue to the metadata Service via SetJobQueue.
func WireJobQueue(svc metadata.Service, queue *Queue) {
	svc.SetJobQueue(queue)
}


