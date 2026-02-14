package sonarr

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
)

// Client is a client for the Sonarr API v3.
// Sonarr is a PRIMARY metadata provider - local, no proxy needed.
type Client struct {
	client      *req.Client
	baseURL     string
	apiKey      string
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
}

// Config contains configuration for the Sonarr client.
type Config struct {
	BaseURL   string
	APIKey    string
	RateLimit rate.Limit // requests per second
	CacheTTL  time.Duration
	Timeout   time.Duration
}

// NewClient creates a new Sonarr API client.
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

	// Circuit breaker (local service — less tolerant)
	circuitbreaker.WrapReqClient(client, "sonarr", circuitbreaker.TierLocal)

	l1, err := cache.NewL1Cache[string, any](1000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		// Fallback: create with defaults if custom config fails.
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
	val, ok := c.cache.Get(key)
	if ok {
		return val
	}
	return nil
}

// setCache stores a value in cache.
func (c *Client) setCache(key string, data any) {
	c.cache.Set(key, data)
}

// GetSystemStatus returns Sonarr system status.
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
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetAllSeries returns all series from Sonarr.
func (c *Client) GetAllSeries(ctx context.Context) ([]Series, error) {
	cacheKey := "series:all"
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Series); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetSeries returns a specific series by ID.
func (c *Client) GetSeries(ctx context.Context, seriesID int) (*Series, error) {
	cacheKey := fmt.Sprintf("series:%d", seriesID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Series); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/series/%d", seriesID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrSeriesNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetSeriesByTVDbID returns a series by TVDb ID.
func (c *Client) GetSeriesByTVDbID(ctx context.Context, tvdbID int) (*Series, error) {
	cacheKey := fmt.Sprintf("series:tvdb:%d", tvdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Series); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("tvdbId", fmt.Sprintf("%d", tvdbID)).
		SetSuccessResult(&result).
		Get("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	if len(result) == 0 {
		return nil, ErrSeriesNotFound
	}

	c.setCache(cacheKey, &result[0])
	return &result[0], nil
}

// GetEpisodes returns all episodes for a series.
func (c *Client) GetEpisodes(ctx context.Context, seriesID int) ([]Episode, error) {
	cacheKey := fmt.Sprintf("series:%d:episodes", seriesID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Episode); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Episode
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("seriesId", fmt.Sprintf("%d", seriesID)).
		SetSuccessResult(&result).
		Get("/episode")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetEpisode returns a specific episode by ID.
func (c *Client) GetEpisode(ctx context.Context, episodeID int) (*Episode, error) {
	cacheKey := fmt.Sprintf("episode:%d", episodeID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Episode); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Episode
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/episode/%d", episodeID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrEpisodeNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetEpisodesBySeason returns episodes for a specific season.
func (c *Client) GetEpisodesBySeason(ctx context.Context, seriesID, seasonNumber int) ([]Episode, error) {
	cacheKey := fmt.Sprintf("series:%d:season:%d:episodes", seriesID, seasonNumber)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Episode); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Episode
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"seriesId":     fmt.Sprintf("%d", seriesID),
			"seasonNumber": fmt.Sprintf("%d", seasonNumber),
		}).
		SetSuccessResult(&result).
		Get("/episode")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetEpisodeFiles returns all episode files for a series.
func (c *Client) GetEpisodeFiles(ctx context.Context, seriesID int) ([]EpisodeFile, error) {
	cacheKey := fmt.Sprintf("series:%d:files", seriesID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]EpisodeFile); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []EpisodeFile
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("seriesId", fmt.Sprintf("%d", seriesID)).
		SetSuccessResult(&result).
		Get("/episodefile")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetEpisodeFile returns a specific episode file by ID.
func (c *Client) GetEpisodeFile(ctx context.Context, episodeFileID int) (*EpisodeFile, error) {
	cacheKey := fmt.Sprintf("episodefile:%d", episodeFileID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*EpisodeFile); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result EpisodeFile
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/episodefile/%d", episodeFileID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, ErrEpisodeFileNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
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
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
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
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
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
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetCalendar returns episodes airing in the specified date range.
func (c *Client) GetCalendar(ctx context.Context, start, end time.Time) ([]CalendarEntry, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []CalendarEntry
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"start":                start.Format(time.RFC3339),
			"end":                  end.Format(time.RFC3339),
			"unmonitored":          "false",
			"includeSeries":        "true",
			"includeEpisodeFile":   "false",
			"includeEpisodeImages": "false",
		}).
		SetSuccessResult(&result).
		Get("/calendar")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	return result, nil
}

// GetHistory returns episode history with pagination.
func (c *Client) GetHistory(ctx context.Context, page, pageSize int, seriesID, episodeID *int) (*HistoryResponse, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	params := map[string]string{
		"page":          fmt.Sprintf("%d", page),
		"pageSize":      fmt.Sprintf("%d", pageSize),
		"sortKey":       "date",
		"sortDirection": "descending",
	}
	if seriesID != nil {
		params["seriesId"] = fmt.Sprintf("%d", *seriesID)
	}
	if episodeID != nil {
		params["episodeId"] = fmt.Sprintf("%d", *episodeID)
	}

	var result HistoryResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetSuccessResult(&result).
		Get("/history")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	return &result, nil
}

// AddSeries adds a new series to Sonarr.
func (c *Client) AddSeries(ctx context.Context, req AddSeriesRequest) (*Series, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetSuccessResult(&result).
		Post("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status, resp.String())
	}

	// Invalidate cache
	c.cache.Delete("series:all")

	return &result, nil
}

// DeleteSeries deletes a series from Sonarr.
func (c *Client) DeleteSeries(ctx context.Context, seriesID int, deleteFiles, addImportListExclusion bool) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"deleteFiles":            fmt.Sprintf("%t", deleteFiles),
			"addImportListExclusion": fmt.Sprintf("%t", addImportListExclusion),
		}).
		Delete(fmt.Sprintf("/series/%d", seriesID))

	if err != nil {
		return fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return ErrSeriesNotFound
		}
		return fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	// Invalidate cache
	c.cache.Delete("series:all")
	c.cache.Delete(fmt.Sprintf("series:%d", seriesID))

	return nil
}

// RefreshSeries triggers a metadata refresh for a series.
func (c *Client) RefreshSeries(ctx context.Context, seriesID int) (*Command, error) {
	return c.runCommand(ctx, "RefreshSeries", &seriesID, nil)
}

// RescanSeries triggers a file rescan for a series.
func (c *Client) RescanSeries(ctx context.Context, seriesID int) (*Command, error) {
	return c.runCommand(ctx, "RescanSeries", &seriesID, nil)
}

// SearchSeries triggers a search for all episodes of a series.
func (c *Client) SearchSeries(ctx context.Context, seriesID int) (*Command, error) {
	return c.runCommand(ctx, "SeriesSearch", &seriesID, nil)
}

// SearchSeason triggers a search for all episodes of a season.
func (c *Client) SearchSeason(ctx context.Context, seriesID, seasonNumber int) (*Command, error) {
	return c.runCommandWithSeason(ctx, "SeasonSearch", seriesID, seasonNumber)
}

// SearchEpisodes triggers a search for specific episodes.
func (c *Client) SearchEpisodes(ctx context.Context, episodeIDs []int) (*Command, error) {
	return c.runCommand(ctx, "EpisodeSearch", nil, episodeIDs)
}

// runCommand executes a command in Sonarr.
func (c *Client) runCommand(ctx context.Context, name string, seriesID *int, episodeIDs []int) (*Command, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	body := map[string]any{
		"name": name,
	}
	if seriesID != nil {
		body["seriesId"] = *seriesID
	}
	if len(episodeIDs) > 0 {
		body["episodeIds"] = episodeIDs
	}

	var result Command
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetSuccessResult(&result).
		Post("/command")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status, resp.String())
	}

	return &result, nil
}

// runCommandWithSeason executes a command with season parameter.
func (c *Client) runCommandWithSeason(ctx context.Context, name string, seriesID, seasonNumber int) (*Command, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	body := map[string]any{
		"name":         name,
		"seriesId":     seriesID,
		"seasonNumber": seasonNumber,
	}

	var result Command
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetSuccessResult(&result).
		Post("/command")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status, resp.String())
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
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status)
	}

	return &result, nil
}

// LookupSeries searches for TV series via Sonarr's lookup API.
// This calls the user's Sonarr instance which internally uses its metadata sources.
func (c *Client) LookupSeries(ctx context.Context, term string) ([]Series, error) {
	cacheKey := fmt.Sprintf("lookup:term:%s", term)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Series); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result []Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("term", term).
		SetSuccessResult(&result).
		Get("/series/lookup")

	if err != nil {
		return nil, fmt.Errorf("sonarr lookup request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// LookupSeriesByTVDbID looks up a series in Sonarr's metadata by TVDb ID.
func (c *Client) LookupSeriesByTVDbID(ctx context.Context, tvdbID int) (*Series, error) {
	cacheKey := fmt.Sprintf("lookup:tvdb:%d", tvdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Series); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	// Sonarr lookup by tvdbid is done via the term parameter with prefix
	var results []Series
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("term", fmt.Sprintf("tvdb:%d", tvdbID)).
		SetSuccessResult(&results).
		Get("/series/lookup")

	if err != nil {
		return nil, fmt.Errorf("sonarr lookup request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("sonarr lookup error: %s", resp.Status)
	}

	if len(results) == 0 {
		return nil, ErrSeriesNotFound
	}

	c.setCache(cacheKey, &results[0])
	return &results[0], nil
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

// IsHealthy checks if Sonarr is reachable and healthy.
func (c *Client) IsHealthy(ctx context.Context) bool {
	status, err := c.GetSystemStatus(ctx)
	if err != nil {
		return false
	}
	return status.Version != ""
}
