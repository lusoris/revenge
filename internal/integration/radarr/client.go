package radarr

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// Client is a client for the Radarr API v3.
// Radarr is a PRIMARY metadata provider - local, no proxy needed.
type Client struct {
	client      *req.Client
	baseURL     string
	apiKey      string
	rateLimiter *rate.Limiter
	cache *cache.L1Cache[string, any]
}

// Config contains configuration for the Radarr client.
type Config struct {
	BaseURL   string
	APIKey    string
	RateLimit rate.Limit // requests per second
	CacheTTL  time.Duration
	Timeout   time.Duration
}

// NewClient creates a new Radarr API client.
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(10.0) // 10 req/s default for local service
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 5 * time.Minute // Short TTL for local cache
	}

	client := req.C().
		SetBaseURL(config.BaseURL+"/api/v3").
		SetTimeout(config.Timeout).
		SetCommonHeader("X-Api-Key", config.APIKey).
		SetCommonHeader("Content-Type", "application/json").
		SetCommonRetryCount(3).
		SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	l1, err := cache.NewL1Cache[string, any](1000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		// Fallback: create with defaults if configuration fails.
		l1, _ = cache.NewL1Cache[string, any](0, 0)
	}

	return &Client{
		client:      client,
		baseURL:     config.BaseURL,
		apiKey:      config.APIKey,
		rateLimiter: rate.NewLimiter(config.RateLimit, 20), // burst of 20
		cache:       l1,
	}
}

// getFromCache retrieves a value from cache.
func (c *Client) getFromCache(key string) any {
	if val, ok := c.cache.Get(key); ok {
		return val
	}
	return nil
}

// setCache stores a value in cache.
func (c *Client) setCache(key string, data any) {
	c.cache.Set(key, data)
}

// GetSystemStatus returns Radarr system status.
func (c *Client) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	cacheKey := "system:status"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*SystemStatus); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result SystemStatus
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/system/status")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetAllMovies returns all movies from Radarr.
func (c *Client) GetAllMovies(ctx context.Context) ([]Movie, error) {
	cacheKey := "movies:all"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/movie")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetMovie returns a specific movie by ID.
func (c *Client) GetMovie(ctx context.Context, movieID int) (*Movie, error) {
	cacheKey := fmt.Sprintf("movie:%d", movieID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movie/%d", movieID))

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetMovieByTMDbID returns a movie by TMDb ID.
func (c *Client) GetMovieByTMDbID(ctx context.Context, tmdbID int) (*Movie, error) {
	cacheKey := fmt.Sprintf("movie:tmdb:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("tmdbId", fmt.Sprintf("%d", tmdbID)).
		SetSuccessResult(&result).
		Get("/movie")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	if len(result) == 0 {
		return nil, ErrMovieNotFound
	}

	c.setCache(cacheKey, &result[0])
	return &result[0], nil
}

// GetMovieFiles returns all movie files for a movie.
func (c *Client) GetMovieFiles(ctx context.Context, movieID int) ([]MovieFile, error) {
	cacheKey := fmt.Sprintf("movie:%d:files", movieID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]MovieFile); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []MovieFile
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("movieId", fmt.Sprintf("%d", movieID)).
		SetSuccessResult(&result).
		Get("/moviefile")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetQualityProfiles returns all quality profiles.
func (c *Client) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	cacheKey := "qualityprofiles"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]QualityProfile); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []QualityProfile
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/qualityprofile")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetRootFolders returns all root folders.
func (c *Client) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	cacheKey := "rootfolders"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]RootFolder); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []RootFolder
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/rootfolder")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetTags returns all tags.
func (c *Client) GetTags(ctx context.Context) ([]Tag, error) {
	cacheKey := "tags"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Tag); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Tag
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/tag")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetCalendar returns movies with releases in the specified date range.
func (c *Client) GetCalendar(ctx context.Context, start, end time.Time) ([]CalendarEntry, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []CalendarEntry
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"start":              start.Format(time.RFC3339),
			"end":                end.Format(time.RFC3339),
			"unmonitored":        "false",
			"includeSeries":      "false",
			"includeEpisodeFile": "false",
			"includeEpisodeImages": "false",
		}).
		SetSuccessResult(&result).
		Get("/calendar")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	return result, nil
}

// GetHistory returns movie history with pagination.
func (c *Client) GetHistory(ctx context.Context, page, pageSize int, movieID *int) (*HistoryResponse, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	params := map[string]string{
		"page":          fmt.Sprintf("%d", page),
		"pageSize":      fmt.Sprintf("%d", pageSize),
		"sortKey":       "date",
		"sortDirection": "descending",
	}
	if movieID != nil {
		params["movieId"] = fmt.Sprintf("%d", *movieID)
	}

	var result HistoryResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/history")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	return &result, nil
}

// AddMovie adds a new movie to Radarr.
func (c *Client) AddMovie(ctx context.Context, req AddMovieRequest) (*Movie, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetSuccessResult(&result).
		Post("/movie")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s - %s", resp.Status, resp.String())
	}

	// Invalidate cache
	c.cache.Delete("movies:all")

	return &result, nil
}

// DeleteMovie deletes a movie from Radarr.
func (c *Client) DeleteMovie(ctx context.Context, movieID int, deleteFiles, addImportExclusion bool) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"deleteFiles":        fmt.Sprintf("%t", deleteFiles),
			"addImportExclusion": fmt.Sprintf("%t", addImportExclusion),
		}).
		Delete(fmt.Sprintf("/movie/%d", movieID))

	if err != nil {
		return fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return ErrMovieNotFound
		}
		return fmt.Errorf("radarr api error: %s", resp.Status)
	}

	// Invalidate cache
	c.cache.Delete("movies:all")
	c.cache.Delete(fmt.Sprintf("movie:%d", movieID))

	return nil
}

// RefreshMovie triggers a metadata refresh for a movie.
func (c *Client) RefreshMovie(ctx context.Context, movieID int) (*Command, error) {
	return c.runCommand(ctx, "RefreshMovie", &movieID)
}

// RescanMovie triggers a file rescan for a movie.
func (c *Client) RescanMovie(ctx context.Context, movieID int) (*Command, error) {
	return c.runCommand(ctx, "RescanMovie", &movieID)
}

// SearchMovie triggers a search for a movie.
func (c *Client) SearchMovie(ctx context.Context, movieID int) (*Command, error) {
	return c.runCommand(ctx, "MoviesSearch", &movieID)
}

// runCommand executes a command in Radarr.
func (c *Client) runCommand(ctx context.Context, name string, movieID *int) (*Command, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	body := map[string]any{
		"name": name,
	}
	if movieID != nil {
		body["movieIds"] = []int{*movieID}
	}

	var result Command
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetSuccessResult(&result).
		Post("/command")

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s - %s", resp.Status, resp.String())
	}

	return &result, nil
}

// GetCommand gets the status of a command.
func (c *Client) GetCommand(ctx context.Context, commandID int) (*Command, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Command
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/command/%d", commandID))

	if err != nil {
		return nil, fmt.Errorf("radarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr api error: %s", resp.Status)
	}

	return &result, nil
}

// LookupMovie searches for movies via Radarr's lookup API.
// This calls the user's Radarr instance which internally uses its metadata sources.
func (c *Client) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	cacheKey := fmt.Sprintf("lookup:term:%s", term)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("term", term).
		SetSuccessResult(&result).
		Get("/movie/lookup")

	if err != nil {
		return nil, fmt.Errorf("radarr lookup request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("radarr lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// LookupMovieByTMDbID looks up a movie in Radarr's metadata by TMDb ID.
func (c *Client) LookupMovieByTMDbID(ctx context.Context, tmdbID int) (*Movie, error) {
	cacheKey := fmt.Sprintf("lookup:tmdb:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movie/lookup/tmdb/%d", tmdbID))

	if err != nil {
		return nil, fmt.Errorf("radarr lookup request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("radarr lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// LookupMovieByIMDbID looks up a movie in Radarr's metadata by IMDb ID.
func (c *Client) LookupMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	cacheKey := fmt.Sprintf("lookup:imdb:%s", imdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movie/lookup/imdb/%s", imdbID))

	if err != nil {
		return nil, fmt.Errorf("radarr lookup request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("radarr lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// ClearCache clears all cached data.
func (c *Client) ClearCache() {
	c.cache.Clear()
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}

// IsHealthy checks if Radarr is reachable and healthy.
func (c *Client) IsHealthy(ctx context.Context) bool {
	status, err := c.GetSystemStatus(ctx)
	if err != nil {
		return false
	}
	return status.Version != ""
}
