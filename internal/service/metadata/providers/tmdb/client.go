package tmdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

const (
	// BaseURL is the TMDb API base URL.
	BaseURL = "https://api.themoviedb.org/3"

	// ImageBaseURL is the TMDb image CDN base URL.
	ImageBaseURL = "https://image.tmdb.org/t/p"

	// DefaultRateLimit is 40 requests per 10 seconds.
	DefaultRateLimit = rate.Limit(4.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 10

	// DefaultCacheTTL is the default cache duration.
	DefaultCacheTTL = 24 * time.Hour

	// SearchCacheTTL is the cache duration for search results.
	SearchCacheTTL = 15 * time.Minute
)

// Config configures the TMDb client.
type Config struct {
	// APIKey is the TMDb API key (v3).
	APIKey string

	// AccessToken is the TMDb access token (v4, optional).
	AccessToken string

	// RateLimit is requests per second (default: 4).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 10).
	Burst int

	// CacheTTL is the cache duration (default: 24h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 30s).
	Timeout time.Duration

	// ProxyURL for HTTP proxy (optional).
	ProxyURL string

	// RetryCount is the number of retries (default: 3).
	RetryCount int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		RateLimit:  DefaultRateLimit,
		Burst:      DefaultBurst,
		CacheTTL:   DefaultCacheTTL,
		Timeout:    30 * time.Second,
		RetryCount: 3,
	}
}

// Client is the TMDb API client with rate limiting and caching.
type Client struct {
	httpClient  *resty.Client
	apiKey      string
	accessToken string
	rateLimiter *rate.Limiter
	cache       sync.Map
	cacheTTL    time.Duration
}

// NewClient creates a new TMDb client.
func NewClient(config Config) *Client {
	// Apply defaults
	if config.RateLimit == 0 {
		config.RateLimit = DefaultRateLimit
	}
	if config.Burst == 0 {
		config.Burst = DefaultBurst
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = DefaultCacheTTL
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	client := resty.New().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	if config.ProxyURL != "" {
		client.SetProxy(config.ProxyURL)
	}

	return &Client{
		httpClient:  client,
		apiKey:      config.APIKey,
		accessToken: config.AccessToken,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cacheTTL:    config.CacheTTL,
	}
}

// request creates an authenticated request.
func (c *Client) request(ctx context.Context) *resty.Request {
	req := c.httpClient.R().SetContext(ctx)

	// Use access token if available, otherwise API key
	if c.accessToken != "" {
		req.SetAuthToken(c.accessToken)
	} else {
		req.SetQueryParam("api_key", c.apiKey)
	}

	return req
}

// waitRateLimit waits for the rate limiter.
func (c *Client) waitRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

// getFromCache retrieves a value from cache.
func (c *Client) getFromCache(key string) any {
	if val, ok := c.cache.Load(key); ok {
		entry, ok := val.(*CacheEntry)
		if !ok {
			return nil
		}
		if !entry.IsExpired() {
			return entry.Data
		}
		c.cache.Delete(key)
	}
	return nil
}

// setCache stores a value in cache.
func (c *Client) setCache(key string, data any, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.cacheTTL
	}
	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.cache.Store(key, entry)
}

// clearCache removes all cached entries.
func (c *Client) clearCache() {
	c.cache.Range(func(key, _ any) bool {
		c.cache.Delete(key)
		return true
	})
}

// cacheKey generates a cache key from components.
func cacheKey(parts ...any) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", part)
	}
	return key
}

// parseError converts API response to error.
func (c *Client) parseError(resp *resty.Response, errResp *ErrorResponse) error {
	if errResp != nil && errResp.StatusMessage != "" {
		return fmt.Errorf("tmdb api error %d: %s", resp.StatusCode(), errResp.StatusMessage)
	}
	return fmt.Errorf("tmdb api error: status %d", resp.StatusCode())
}

// SearchMovie searches for movies.
func (c *Client) SearchMovie(ctx context.Context, query string, year *int, language string) (*SearchResultsResponse, error) {
	key := cacheKey("search:movie", query, year, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SearchResultsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := map[string]string{
		"query": query,
	}
	if year != nil {
		params["year"] = fmt.Sprintf("%d", *year)
	}
	if language != "" {
		params["language"] = language
	}

	var result SearchResultsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get("/search/movie")

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// GetMovie retrieves movie details.
func (c *Client) GetMovie(ctx context.Context, id int, language string, appendToResponse string) (*MovieResponse, error) {
	key := cacheKey("movie", id, language, appendToResponse)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*MovieResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}
	if appendToResponse != "" {
		params["append_to_response"] = appendToResponse
	}

	var result MovieResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetMovieCredits retrieves movie credits.
func (c *Client) GetMovieCredits(ctx context.Context, id int) (*CreditsResponse, error) {
	key := cacheKey("movie:credits", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*CreditsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result CreditsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/credits", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetMovieImages retrieves movie images.
func (c *Client) GetMovieImages(ctx context.Context, id int) (*ImagesResponse, error) {
	key := cacheKey("movie:images", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ImagesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ImagesResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/images", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetMovieReleaseDates retrieves movie release dates.
func (c *Client) GetMovieReleaseDates(ctx context.Context, id int) (*ReleaseDatesWrapper, error) {
	key := cacheKey("movie:releases", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ReleaseDatesWrapper); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ReleaseDatesWrapper
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/release_dates", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetMovieTranslations retrieves movie translations.
func (c *Client) GetMovieTranslations(ctx context.Context, id int) (*TranslationsWrapper, error) {
	key := cacheKey("movie:translations", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*TranslationsWrapper); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result TranslationsWrapper
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/translations", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetMovieExternalIDs retrieves movie external IDs.
func (c *Client) GetMovieExternalIDs(ctx context.Context, id int) (*ExternalIDsResponse, error) {
	key := cacheKey("movie:external", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ExternalIDsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ExternalIDsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/external_ids", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetSimilarMovies retrieves movies similar to the given movie.
func (c *Client) GetSimilarMovies(ctx context.Context, id int, language string, page int) (*SearchResultsResponse, error) {
	if page < 1 {
		page = 1
	}
	key := cacheKey("movie:similar", id, language, page)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SearchResultsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
	}
	if language != "" {
		params["language"] = language
	}

	var result SearchResultsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/similar", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// GetMovieRecommendations retrieves recommended movies based on the given movie.
func (c *Client) GetMovieRecommendations(ctx context.Context, id int, language string, page int) (*SearchResultsResponse, error) {
	if page < 1 {
		page = 1
	}
	key := cacheKey("movie:recommendations", id, language, page)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SearchResultsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
	}
	if language != "" {
		params["language"] = language
	}

	var result SearchResultsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/recommendations", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// SearchTV searches for TV shows.
func (c *Client) SearchTV(ctx context.Context, query string, year *int, language string) (*TVSearchResultsResponse, error) {
	key := cacheKey("search:tv", query, year, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*TVSearchResultsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := map[string]string{
		"query": query,
	}
	if year != nil {
		params["first_air_date_year"] = fmt.Sprintf("%d", *year)
	}
	if language != "" {
		params["language"] = language
	}

	var result TVSearchResultsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get("/search/tv")

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// GetTV retrieves TV show details.
func (c *Client) GetTV(ctx context.Context, id int, language string, appendToResponse string) (*TVResponse, error) {
	key := cacheKey("tv", id, language, appendToResponse)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*TVResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}
	if appendToResponse != "" {
		params["append_to_response"] = appendToResponse
	}

	var result TVResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetTVCredits retrieves TV show credits.
func (c *Client) GetTVCredits(ctx context.Context, id int) (*CreditsResponse, error) {
	key := cacheKey("tv:credits", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*CreditsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result CreditsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/credits", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetTVImages retrieves TV show images.
func (c *Client) GetTVImages(ctx context.Context, id int) (*ImagesResponse, error) {
	key := cacheKey("tv:images", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ImagesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ImagesResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/images", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetTVContentRatings retrieves TV show content ratings.
func (c *Client) GetTVContentRatings(ctx context.Context, id int) (*ContentRatingsWrapper, error) {
	key := cacheKey("tv:ratings", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ContentRatingsWrapper); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ContentRatingsWrapper
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/content_ratings", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetTVTranslations retrieves TV show translations.
func (c *Client) GetTVTranslations(ctx context.Context, id int) (*TranslationsWrapper, error) {
	key := cacheKey("tv:translations", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*TranslationsWrapper); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result TranslationsWrapper
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/translations", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetTVExternalIDs retrieves TV show external IDs.
func (c *Client) GetTVExternalIDs(ctx context.Context, id int) (*ExternalIDsResponse, error) {
	key := cacheKey("tv:external", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ExternalIDsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ExternalIDsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/external_ids", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetSeason retrieves season details.
func (c *Client) GetSeason(ctx context.Context, tvID, seasonNum int, language string, appendToResponse string) (*SeasonResponse, error) {
	key := cacheKey("tv:season", tvID, seasonNum, language, appendToResponse)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SeasonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}
	if appendToResponse != "" {
		params["append_to_response"] = appendToResponse
	}

	var result SeasonResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d", tvID, seasonNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetSeasonCredits retrieves season credits.
func (c *Client) GetSeasonCredits(ctx context.Context, tvID, seasonNum int) (*CreditsResponse, error) {
	key := cacheKey("tv:season:credits", tvID, seasonNum)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*CreditsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result CreditsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d/credits", tvID, seasonNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetSeasonImages retrieves season images.
func (c *Client) GetSeasonImages(ctx context.Context, tvID, seasonNum int) (*ImagesResponse, error) {
	key := cacheKey("tv:season:images", tvID, seasonNum)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ImagesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ImagesResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d/images", tvID, seasonNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetEpisode retrieves episode details.
func (c *Client) GetEpisode(ctx context.Context, tvID, seasonNum, episodeNum int, language string, appendToResponse string) (*EpisodeResponse, error) {
	key := cacheKey("tv:episode", tvID, seasonNum, episodeNum, language, appendToResponse)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*EpisodeResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}
	if appendToResponse != "" {
		params["append_to_response"] = appendToResponse
	}

	var result EpisodeResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d/episode/%d", tvID, seasonNum, episodeNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetEpisodeCredits retrieves episode credits.
func (c *Client) GetEpisodeCredits(ctx context.Context, tvID, seasonNum, episodeNum int) (*CreditsResponse, error) {
	key := cacheKey("tv:episode:credits", tvID, seasonNum, episodeNum)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*CreditsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result CreditsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d/episode/%d/credits", tvID, seasonNum, episodeNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetEpisodeImages retrieves episode images.
func (c *Client) GetEpisodeImages(ctx context.Context, tvID, seasonNum, episodeNum int) (*ImagesResponse, error) {
	key := cacheKey("tv:episode:images", tvID, seasonNum, episodeNum)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ImagesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ImagesResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/tv/%d/season/%d/episode/%d/images", tvID, seasonNum, episodeNum))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// SearchPerson searches for people.
func (c *Client) SearchPerson(ctx context.Context, query string, language string) (*PersonSearchResultsResponse, error) {
	key := cacheKey("search:person", query, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonSearchResultsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := map[string]string{
		"query": query,
	}
	if language != "" {
		params["language"] = language
	}

	var result PersonSearchResultsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get("/search/person")

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// GetPerson retrieves person details.
func (c *Client) GetPerson(ctx context.Context, id int, language string, appendToResponse string) (*PersonResponse, error) {
	key := cacheKey("person", id, language, appendToResponse)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}
	if appendToResponse != "" {
		params["append_to_response"] = appendToResponse
	}

	var result PersonResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/person/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetPersonCredits retrieves person credits (filmography).
func (c *Client) GetPersonCredits(ctx context.Context, id int, language string) (*PersonCreditsResponse, error) {
	key := cacheKey("person:credits", id, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonCreditsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}

	var result PersonCreditsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/person/%d/combined_credits", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetPersonImages retrieves person images.
func (c *Client) GetPersonImages(ctx context.Context, id int) (*PersonImagesResponse, error) {
	key := cacheKey("person:images", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonImagesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result PersonImagesResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/person/%d/images", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetPersonExternalIDs retrieves person external IDs.
func (c *Client) GetPersonExternalIDs(ctx context.Context, id int) (*ExternalIDsResponse, error) {
	key := cacheKey("person:external", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*ExternalIDsResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result ExternalIDsResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/person/%d/external_ids", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetCollection retrieves collection details.
func (c *Client) GetCollection(ctx context.Context, id int, language string) (*CollectionResponse, error) {
	key := cacheKey("collection", id, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*CollectionResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	params := make(map[string]string)
	if language != "" {
		params["language"] = language
	}

	var result CollectionResponse
	var errResp ErrorResponse

	resp, err := c.request(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/collection/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, 0)
	return &result, nil
}

// GetImageURL constructs a full image URL.
func (c *Client) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s%s", ImageBaseURL, size, path)
}

// DownloadImage downloads an image.
func (c *Client) DownloadImage(ctx context.Context, path string, size string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("empty image path")
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	url := c.GetImageURL(path, size)

	// Use separate client for image CDN
	resp, err := resty.New().R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("download image: status %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// ClearCache clears all cached data.
func (c *Client) ClearCache() {
	c.clearCache()
}
