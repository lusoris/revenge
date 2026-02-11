package jobs

import (
	"context"

	"github.com/google/uuid"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
)

// Queue provides an interface for enqueuing metadata jobs.
type Queue struct {
	client *infrajobs.Client
}

// NewQueue creates a new metadata job queue.
func NewQueue(client *infrajobs.Client) *Queue {
	return &Queue{client: client}
}

// EnqueueRefreshMovie enqueues a movie metadata refresh job.
func (q *Queue) EnqueueRefreshMovie(ctx context.Context, movieID uuid.UUID, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshMovieArgs{
		MovieID:   movieID,
		Force:     force,
		Languages: languages,
	}, nil)
	return err
}

// EnqueueRefreshTVShow enqueues a TV show metadata refresh job.
func (q *Queue) EnqueueRefreshTVShow(ctx context.Context, seriesID uuid.UUID, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshTVShowArgs{
		SeriesID:        seriesID,
		Force:           force,
		Languages:       languages,
		IncludeSeasons:  true,
		IncludeEpisodes: false,
	}, nil)
	return err
}

// EnqueueRefreshTVShowFull enqueues a full TV show metadata refresh job including episodes.
func (q *Queue) EnqueueRefreshTVShowFull(ctx context.Context, seriesID uuid.UUID, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshTVShowArgs{
		SeriesID:        seriesID,
		Force:           force,
		Languages:       languages,
		IncludeSeasons:  true,
		IncludeEpisodes: true,
	}, nil)
	return err
}

// EnqueueRefreshSeason enqueues a season metadata refresh job.
func (q *Queue) EnqueueRefreshSeason(ctx context.Context, seriesID, seasonID uuid.UUID, seasonNum int, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshSeasonArgs{
		SeriesID:        seriesID,
		SeasonID:        seasonID,
		SeasonNumber:    seasonNum,
		Force:           force,
		Languages:       languages,
		IncludeEpisodes: true,
	}, nil)
	return err
}

// EnqueueRefreshEpisode enqueues an episode metadata refresh job.
func (q *Queue) EnqueueRefreshEpisode(ctx context.Context, seriesID, seasonID, episodeID uuid.UUID, seasonNum, episodeNum int, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshEpisodeArgs{
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		EpisodeID:     episodeID,
		SeasonNumber:  seasonNum,
		EpisodeNumber: episodeNum,
		Force:         force,
		Languages:     languages,
	}, nil)
	return err
}

// EnqueueRefreshPerson enqueues a person metadata refresh job.
func (q *Queue) EnqueueRefreshPerson(ctx context.Context, personID uuid.UUID, providerID string, force bool, languages []string) error {
	_, err := q.client.Insert(ctx, RefreshPersonArgs{
		PersonID:   personID,
		ProviderID: providerID,
		Force:      force,
		Languages:  languages,
	}, nil)
	return err
}

// EnqueueEnrichContent enqueues a content enrichment job.
func (q *Queue) EnqueueEnrichContent(ctx context.Context, contentType string, contentID uuid.UUID, providers, languages []string) error {
	_, err := q.client.Insert(ctx, EnrichContentArgs{
		ContentType: contentType,
		ContentID:   contentID,
		Providers:   providers,
		Languages:   languages,
	}, nil)
	return err
}

// EnqueueDownloadImage enqueues an image download job.
func (q *Queue) EnqueueDownloadImage(ctx context.Context, contentType, contentID, imageType, path, size string) error {
	_, err := q.client.Insert(ctx, DownloadImageArgs{
		ContentType: contentType,
		ContentID:   contentID,
		ImageType:   imageType,
		Path:        path,
		Size:        size,
	}, nil)
	return err
}

// BatchEnqueueRefreshMovies enqueues multiple movie refresh jobs.
func (q *Queue) BatchEnqueueRefreshMovies(ctx context.Context, movieIDs []uuid.UUID, force bool, languages []string) error {
	params := make([]river.InsertManyParams, len(movieIDs))
	for i, id := range movieIDs {
		params[i] = river.InsertManyParams{
			Args: RefreshMovieArgs{
				MovieID:   id,
				Force:     force,
				Languages: languages,
			},
		}
	}
	_, err := q.client.InsertMany(ctx, params)
	return err
}

// BatchEnqueueRefreshTVShows enqueues multiple TV show refresh jobs.
func (q *Queue) BatchEnqueueRefreshTVShows(ctx context.Context, seriesIDs []uuid.UUID, force bool, languages []string) error {
	params := make([]river.InsertManyParams, len(seriesIDs))
	for i, id := range seriesIDs {
		params[i] = river.InsertManyParams{
			Args: RefreshTVShowArgs{
				SeriesID:        id,
				Force:           force,
				Languages:       languages,
				IncludeSeasons:  true,
				IncludeEpisodes: false,
			},
		}
	}
	_, err := q.client.InsertMany(ctx, params)
	return err
}
