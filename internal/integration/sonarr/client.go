package sonarr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

// Client is a client for the Sonarr API v3.
// Sonarr is a PRIMARY metadata provider - local, no proxy needed.
type Client struct {
	client      *resty.Client
	baseURL     string
	apiKey      string
	rateLimiter *rate.Limiter
	cache       sync.Map
	cacheTTL    time.Duration
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

	client := resty.New().
		SetBaseURL(config.BaseURL+"/api/v3").
		SetTimeout(config.Timeout).
		SetHeader("X-Api-Key", config.APIKey).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	return &Client{
		client:      client,
		baseURL:     config.BaseURL,
		apiKey:      config.APIKey,
		rateLimiter: rate.NewLimiter(config.RateLimit, 20), // burst of 20
		cacheTTL:    config.CacheTTL,
	}
}

// cacheEntry represents a cached item with expiration.
type cacheEntry struct {
	data      any
	expiresAt time.Time
}

// getFromCache retrieves a value from cache if not expired.
func (c *Client) getFromCache(key string) any {
	if entry, ok := c.cache.Load(key); ok {
		e, ok := entry.(cacheEntry)
		if !ok {
			return nil
		}
		if time.Now().Before(e.expiresAt) {
			return e.data
		}
		c.cache.Delete(key)
	}
	return nil
}

// setCache stores a value in cache.
func (c *Client) setCache(key string, data any) {
	c.cache.Store(key, cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.cacheTTL),
	})
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
		SetResult(&result).
		Get("/system/status")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get(fmt.Sprintf("/series/%d", seriesID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return nil, ErrSeriesNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/episode")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get(fmt.Sprintf("/episode/%d", episodeID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return nil, ErrEpisodeNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/episode")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/episodefile")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get(fmt.Sprintf("/episodefile/%d", episodeFileID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return nil, ErrEpisodeFileNotFound
		}
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/qualityprofile")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/rootfolder")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/tag")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
			"start":              start.Format(time.RFC3339),
			"end":                end.Format(time.RFC3339),
			"unmonitored":        "false",
			"includeSeries":      "true",
			"includeEpisodeFile": "false",
			"includeEpisodeImages": "false",
		}).
		SetResult(&result).
		Get("/calendar")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Get("/history")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Post("/series")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status(), resp.String())
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
			"deleteFiles":           fmt.Sprintf("%t", deleteFiles),
			"addImportListExclusion": fmt.Sprintf("%t", addImportListExclusion),
		}).
		Delete(fmt.Sprintf("/series/%d", seriesID))

	if err != nil {
		return fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return ErrSeriesNotFound
		}
		return fmt.Errorf("sonarr api error: %s", resp.Status())
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
		SetResult(&result).
		Post("/command")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status(), resp.String())
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
		SetResult(&result).
		Post("/command")

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s - %s", resp.Status(), resp.String())
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
		SetResult(&result).
		Get(fmt.Sprintf("/command/%d", commandID))

	if err != nil {
		return nil, fmt.Errorf("sonarr api request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("sonarr api error: %s", resp.Status())
	}

	return &result, nil
}

// ClearCache clears all cached data.
func (c *Client) ClearCache() {
	c.cache = sync.Map{}
}

// IsHealthy checks if Sonarr is reachable and healthy.
func (c *Client) IsHealthy(ctx context.Context) bool {
	status, err := c.GetSystemStatus(ctx)
	if err != nil {
		return false
	}
	return status.Version != ""
}
