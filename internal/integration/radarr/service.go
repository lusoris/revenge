package radarr

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/util"
)

// SyncService handles synchronization between Radarr and Revenge.
// It implements the PRIMARY provider pattern from the metadata priority chain.
type SyncService struct {
	client     *Client
	mapper     *Mapper
	movieRepo  movie.Repository
	logger     *slog.Logger
	syncMu     sync.Mutex
	syncStatus SyncStatus
}

// SyncStatus represents the current sync status.
type SyncStatus struct {
	IsRunning     bool      `json:"is_running"`
	LastSync      time.Time `json:"last_sync,omitempty"`
	LastSyncError string    `json:"last_sync_error,omitempty"`
	MoviesAdded   int       `json:"movies_added"`
	MoviesUpdated int       `json:"movies_updated"`
	MoviesRemoved int       `json:"movies_removed"`
	TotalMovies   int       `json:"total_movies"`
}

// SyncResult contains the result of a sync operation.
type SyncResult struct {
	MoviesAdded   int       `json:"movies_added"`
	MoviesUpdated int       `json:"movies_updated"`
	MoviesRemoved int       `json:"movies_removed"`
	MoviesSkipped int       `json:"movies_skipped"`
	Errors        []string  `json:"errors,omitempty"`
	Duration      time.Duration `json:"duration"`
}

// NewSyncService creates a new Radarr sync service.
func NewSyncService(
	client *Client,
	mapper *Mapper,
	movieRepo movie.Repository,
	logger *slog.Logger,
) *SyncService {
	return &SyncService{
		client:    client,
		mapper:    mapper,
		movieRepo: movieRepo,
		logger:    logger.With("service", "radarr_sync"),
	}
}

// GetStatus returns the current sync status.
func (s *SyncService) GetStatus() SyncStatus {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	return s.syncStatus
}

// IsHealthy checks if Radarr is reachable and healthy.
func (s *SyncService) IsHealthy(ctx context.Context) bool {
	return s.client.IsHealthy(ctx)
}

// GetSystemStatus returns Radarr's system status.
func (s *SyncService) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	return s.client.GetSystemStatus(ctx)
}

// GetQualityProfiles returns all quality profiles from Radarr.
func (s *SyncService) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	return s.client.GetQualityProfiles(ctx)
}

// GetRootFolders returns all root folders from Radarr.
func (s *SyncService) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	return s.client.GetRootFolders(ctx)
}

// LookupMovie searches for movies via Radarr's lookup API.
func (s *SyncService) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	return s.client.LookupMovie(ctx, term)
}

// LookupMovieByTMDbID looks up a movie by TMDb ID via Radarr.
func (s *SyncService) LookupMovieByTMDbID(ctx context.Context, tmdbID int) (*Movie, error) {
	return s.client.LookupMovieByTMDbID(ctx, tmdbID)
}

// LookupMovieByIMDbID looks up a movie by IMDb ID via Radarr.
func (s *SyncService) LookupMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	return s.client.LookupMovieByIMDbID(ctx, imdbID)
}

// SyncLibrary performs a full library sync from Radarr to Revenge.
// This is the main sync operation that imports all movies from Radarr.
func (s *SyncService) SyncLibrary(ctx context.Context) (*SyncResult, error) {
	s.syncMu.Lock()
	if s.syncStatus.IsRunning {
		s.syncMu.Unlock()
		return nil, fmt.Errorf("sync already in progress")
	}
	s.syncStatus.IsRunning = true
	s.syncMu.Unlock()

	defer func() {
		s.syncMu.Lock()
		s.syncStatus.IsRunning = false
		s.syncStatus.LastSync = time.Now()
		s.syncMu.Unlock()
	}()

	start := time.Now()
	s.logger.Info("starting library sync from Radarr")

	result := &SyncResult{}

	// Get all movies from Radarr
	radarrMovies, err := s.client.GetAllMovies(ctx)
	if err != nil {
		s.syncMu.Lock()
		s.syncStatus.LastSyncError = err.Error()
		s.syncMu.Unlock()
		return nil, fmt.Errorf("failed to get movies from Radarr: %w", err)
	}

	s.logger.Info("fetched movies from Radarr", "count", len(radarrMovies))

	// Get existing movies from Revenge by RadarrID
	existingMovies, err := s.getExistingRadarrMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing movies: %w", err)
	}

	// Track which movies we've seen for removal detection
	seenRadarrIDs := make(map[int]bool)

	// Process each movie
	for _, rm := range radarrMovies {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		seenRadarrIDs[rm.ID] = true

		// Only sync movies that have files
		if !rm.HasFile {
			result.MoviesSkipped++
			continue
		}

		if existingID, exists := existingMovies[rm.ID]; exists {
			// Update existing movie
			if err := s.updateMovie(ctx, rm, existingID); err != nil {
				s.logger.Error("failed to update movie", "radarr_id", rm.ID, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("update movie %d: %v", rm.ID, err))
				continue
			}
			result.MoviesUpdated++
		} else {
			// Add new movie
			if err := s.addMovie(ctx, rm); err != nil {
				s.logger.Error("failed to add movie", "radarr_id", rm.ID, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("add movie %d: %v", rm.ID, err))
				continue
			}
			result.MoviesAdded++
		}
	}

	// Find and mark removed movies (movies in Revenge that are no longer in Radarr)
	for radarrID := range existingMovies {
		if !seenRadarrIDs[radarrID] {
			// Movie was removed from Radarr - we could delete or mark as removed
			// For now, we log but don't delete (soft removal)
			s.logger.Info("movie no longer in Radarr", "radarr_id", radarrID)
			result.MoviesRemoved++
		}
	}

	result.Duration = time.Since(start)

	s.syncMu.Lock()
	s.syncStatus.MoviesAdded = result.MoviesAdded
	s.syncStatus.MoviesUpdated = result.MoviesUpdated
	s.syncStatus.MoviesRemoved = result.MoviesRemoved
	s.syncStatus.TotalMovies = len(radarrMovies)
	s.syncStatus.LastSyncError = ""
	s.syncMu.Unlock()

	s.logger.Info("library sync completed",
		"added", result.MoviesAdded,
		"updated", result.MoviesUpdated,
		"removed", result.MoviesRemoved,
		"skipped", result.MoviesSkipped,
		"duration", result.Duration,
	)

	return result, nil
}

// SyncMovie syncs a single movie from Radarr by its Radarr ID.
func (s *SyncService) SyncMovie(ctx context.Context, radarrID int) error {
	rm, err := s.client.GetMovie(ctx, radarrID)
	if err != nil {
		return fmt.Errorf("failed to get movie from Radarr: %w", err)
	}

	// Get existing movie by RadarrID
	existing, err := s.movieRepo.GetMovieByRadarrID(ctx, util.SafeIntToInt32(radarrID))
	if err != nil && err != movie.ErrMovieNotFound {
		return fmt.Errorf("failed to check existing movie: %w", err)
	}

	if existing != nil {
		return s.updateMovie(ctx, *rm, existing.ID)
	}
	return s.addMovie(ctx, *rm)
}

// RefreshMovie triggers a refresh in Radarr for a movie by TMDb ID.
func (s *SyncService) RefreshMovie(ctx context.Context, tmdbID int) (*Command, error) {
	rm, err := s.client.GetMovieByTMDbID(ctx, tmdbID)
	if err != nil {
		return nil, fmt.Errorf("movie not found in Radarr: %w", err)
	}
	return s.client.RefreshMovie(ctx, rm.ID)
}

// getExistingRadarrMovies returns a map of RadarrID -> MovieID for all movies with a RadarrID.
func (s *SyncService) getExistingRadarrMovies(ctx context.Context) (map[int]uuid.UUID, error) {
	// Get all movies with RadarrID set
	// This is a simplified approach - in production you might want pagination
	movies, err := s.movieRepo.ListMovies(ctx, movie.ListFilters{
		Limit:  10000,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[int]uuid.UUID)
	for _, m := range movies {
		if m.RadarrID != nil {
			result[int(*m.RadarrID)] = m.ID
		}
	}
	return result, nil
}

// addMovie creates a new movie from Radarr data.
func (s *SyncService) addMovie(ctx context.Context, rm Movie) error {
	m := s.mapper.ToMovie(&rm)

	// Convert to create params
	createParams := s.movieToCreateParams(m)

	// Create the movie
	created, err := s.movieRepo.CreateMovie(ctx, createParams)
	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	// Add genres
	for _, genreName := range rm.Genres {
		if err := s.movieRepo.AddMovieGenre(ctx, created.ID, 0, genreName); err != nil {
			s.logger.Warn("failed to add genre", "movie_id", created.ID, "genre", genreName, "error", err)
		}
	}

	// Handle collection if present
	if rm.Collection != nil {
		if err := s.syncCollection(ctx, rm.Collection, created.ID); err != nil {
			s.logger.Warn("failed to sync collection", "movie_id", created.ID, "error", err)
		}
	}

	// Sync movie files
	if rm.HasFile {
		if err := s.syncMovieFiles(ctx, rm.ID, created.ID); err != nil {
			s.logger.Warn("failed to sync movie files", "movie_id", created.ID, "error", err)
		}
	}

	s.logger.Debug("added movie from Radarr", "id", created.ID, "title", rm.Title)
	return nil
}

// updateMovie updates an existing movie with Radarr data.
func (s *SyncService) updateMovie(ctx context.Context, rm Movie, existingID uuid.UUID) error {
	m := s.mapper.ToMovie(&rm)
	m.ID = existingID

	updateParams := s.movieToUpdateParams(m)
	if _, err := s.movieRepo.UpdateMovie(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	// Update collection if present
	if rm.Collection != nil {
		if err := s.syncCollection(ctx, rm.Collection, existingID); err != nil {
			s.logger.Warn("failed to sync collection", "movie_id", existingID, "error", err)
		}
	}

	// Sync movie files
	if rm.HasFile {
		if err := s.syncMovieFiles(ctx, rm.ID, existingID); err != nil {
			s.logger.Warn("failed to sync movie files", "movie_id", existingID, "error", err)
		}
	}

	s.logger.Debug("updated movie from Radarr", "id", existingID, "title", rm.Title)
	return nil
}

// syncCollection syncs a collection and links the movie to it.
func (s *SyncService) syncCollection(ctx context.Context, rc *Collection, movieID uuid.UUID) error {
	// Try to find existing collection by TMDb ID
	existing, err := s.movieRepo.GetMovieCollectionByTMDbID(ctx, util.SafeIntToInt32(rc.TMDbID))
	if err != nil && err != movie.ErrCollectionNotFound {
		return err
	}

	var collectionID uuid.UUID

	if existing != nil {
		collectionID = existing.ID
		// Update collection info if needed
		mc := s.mapper.ToMovieCollection(rc)
		updateParams := s.collectionToUpdateParams(mc, collectionID)
		if _, err := s.movieRepo.UpdateMovieCollection(ctx, updateParams); err != nil {
			s.logger.Warn("failed to update collection", "id", collectionID, "error", err)
		}
	} else {
		// Create new collection
		mc := s.mapper.ToMovieCollection(rc)
		createParams := s.collectionToCreateParams(mc)
		created, err := s.movieRepo.CreateMovieCollection(ctx, createParams)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		collectionID = created.ID
	}

	// Link movie to collection
	if err := s.movieRepo.AddMovieToCollection(ctx, collectionID, movieID, nil); err != nil {
		s.logger.Warn("failed to link movie to collection", "movie_id", movieID, "collection_id", collectionID, "error", err)
	}

	return nil
}

// syncMovieFiles syncs movie files from Radarr.
func (s *SyncService) syncMovieFiles(ctx context.Context, radarrMovieID int, movieID uuid.UUID) error {
	files, err := s.client.GetMovieFiles(ctx, radarrMovieID)
	if err != nil {
		return fmt.Errorf("failed to get movie files: %w", err)
	}

	for _, rf := range files {
		mf := s.mapper.ToMovieFile(&rf, movieID)

		// Check if file already exists by RadarrFileID
		existing, err := s.movieRepo.GetMovieFileByRadarrID(ctx, util.SafeIntToInt32(rf.ID))
		if err != nil && err != movie.ErrMovieFileNotFound {
			return err
		}

		if existing != nil {
			// Update existing file
			updateParams := s.movieFileToUpdateParams(mf, existing.ID)
			if _, err := s.movieRepo.UpdateMovieFile(ctx, updateParams); err != nil {
				s.logger.Warn("failed to update movie file", "file_id", existing.ID, "error", err)
			}
		} else {
			// Create new file
			createParams := s.movieFileToCreateParams(mf)
			if _, err := s.movieRepo.CreateMovieFile(ctx, createParams); err != nil {
				s.logger.Warn("failed to create movie file", "movie_id", movieID, "error", err)
			}
		}
	}

	return nil
}

// movieToCreateParams converts a domain movie to create params.
func (s *SyncService) movieToCreateParams(m *movie.Movie) movie.CreateMovieParams {
	params := movie.CreateMovieParams{
		TMDbID:           m.TMDbID,
		IMDbID:           m.IMDbID,
		Title:            m.Title,
		OriginalTitle:    m.OriginalTitle,
		Year:             m.Year,
		Runtime:          m.Runtime,
		Overview:         m.Overview,
		Status:           m.Status,
		OriginalLanguage: m.OriginalLanguage,
		PosterPath:       m.PosterPath,
		BackdropPath:     m.BackdropPath,
		TrailerURL:       m.TrailerURL,
		VoteCount:        m.VoteCount,
		RadarrID:         m.RadarrID,
	}

	if m.ReleaseDate != nil {
		rd := m.ReleaseDate.Format("2006-01-02")
		params.ReleaseDate = &rd
	}
	if m.VoteAverage != nil {
		va := m.VoteAverage.String()
		params.VoteAverage = &va
	}
	if m.Popularity != nil {
		p := m.Popularity.String()
		params.Popularity = &p
	}

	return params
}

// movieToUpdateParams converts a domain movie to update params.
func (s *SyncService) movieToUpdateParams(m *movie.Movie) movie.UpdateMovieParams {
	params := movie.UpdateMovieParams{
		ID:               m.ID,
		TMDbID:           m.TMDbID,
		IMDbID:           m.IMDbID,
		Title:            &m.Title,
		OriginalTitle:    m.OriginalTitle,
		Year:             m.Year,
		Runtime:          m.Runtime,
		Overview:         m.Overview,
		Status:           m.Status,
		OriginalLanguage: m.OriginalLanguage,
		PosterPath:       m.PosterPath,
		BackdropPath:     m.BackdropPath,
		TrailerURL:       m.TrailerURL,
		VoteCount:        m.VoteCount,
		RadarrID:         m.RadarrID,
	}

	if m.ReleaseDate != nil {
		rd := m.ReleaseDate.Format("2006-01-02")
		params.ReleaseDate = &rd
	}
	if m.VoteAverage != nil {
		va := m.VoteAverage.String()
		params.VoteAverage = &va
	}
	if m.Popularity != nil {
		p := m.Popularity.String()
		params.Popularity = &p
	}

	return params
}

// collectionToCreateParams converts a domain collection to create params.
func (s *SyncService) collectionToCreateParams(mc *movie.MovieCollection) movie.CreateMovieCollectionParams {
	return movie.CreateMovieCollectionParams{
		TMDbCollectionID: mc.TMDbCollectionID,
		Name:             mc.Name,
		Overview:         mc.Overview,
		PosterPath:       mc.PosterPath,
		BackdropPath:     mc.BackdropPath,
	}
}

// collectionToUpdateParams converts a domain collection to update params.
func (s *SyncService) collectionToUpdateParams(mc *movie.MovieCollection, id uuid.UUID) movie.UpdateMovieCollectionParams {
	return movie.UpdateMovieCollectionParams{
		ID:               id,
		TMDbCollectionID: mc.TMDbCollectionID,
		Name:             &mc.Name,
		Overview:         mc.Overview,
		PosterPath:       mc.PosterPath,
		BackdropPath:     mc.BackdropPath,
	}
}

// movieFileToCreateParams converts a domain movie file to create params.
func (s *SyncService) movieFileToCreateParams(mf *movie.MovieFile) movie.CreateMovieFileParams {
	return movie.CreateMovieFileParams{
		MovieID:        mf.MovieID,
		FilePath:       mf.FilePath,
		FileSize:       mf.FileSize,
		Resolution:     mf.Resolution,
		QualityProfile: mf.QualityProfile,
		VideoCodec:     mf.VideoCodec,
		AudioCodec:     mf.AudioCodec,
		Container:      mf.Container,
		BitrateKbps:    mf.BitrateKbps,
		RadarrFileID:   mf.RadarrFileID,
	}
}

// movieFileToUpdateParams converts a domain movie file to update params.
func (s *SyncService) movieFileToUpdateParams(mf *movie.MovieFile, id uuid.UUID) movie.UpdateMovieFileParams {
	return movie.UpdateMovieFileParams{
		ID:             id,
		FilePath:       &mf.FilePath,
		FileSize:       &mf.FileSize,
		Resolution:     mf.Resolution,
		QualityProfile: mf.QualityProfile,
		VideoCodec:     mf.VideoCodec,
		AudioCodec:     mf.AudioCodec,
		Container:      mf.Container,
		BitrateKbps:    mf.BitrateKbps,
		RadarrFileID:   mf.RadarrFileID,
	}
}
