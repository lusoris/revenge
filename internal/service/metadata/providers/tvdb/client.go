package tvdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	// BaseURL is the TVDb API v4 base URL.
	BaseURL = "https://api4.thetvdb.com/v4"

	// DefaultRateLimit is more conservative for TVDb.
	DefaultRateLimit = rate.Limit(5.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 10

	// DefaultCacheTTL is the default cache duration.
	DefaultCacheTTL = 24 * time.Hour

	// SearchCacheTTL is the cache duration for search results.
	SearchCacheTTL = 15 * time.Minute

	// TokenRefreshBuffer is how early to refresh the token before expiry.
	TokenRefreshBuffer = 1 * time.Hour
)

// Config configures the TVDb client.
type Config struct {
	// APIKey is the TVDb API key.
	APIKey string

	// PIN is the optional subscriber PIN for additional access.
	PIN string

	// RateLimit is requests per second (default: 5).
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

// Client is the TVDb API client with JWT authentication.
type Client struct {
	httpClient  *req.Client
	apiKey      string
	pin         string
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]

	// JWT token management
	token       string
	tokenExpiry time.Time
	tokenMutex  sync.RWMutex
}

// NewClient creates a new TVDb client.
func NewClient(config Config) (*Client, error) {
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

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("tvdb l1 cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(config.RetryCount).
		SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second)

	if config.ProxyURL != "" {
		client.SetProxyURL(config.ProxyURL)
	}

	return &Client{
		httpClient:  client,
		apiKey:      config.APIKey,
		pin:         config.PIN,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cache:       l1,
	}, nil
}

// authenticate obtains a JWT token from TVDb.
func (c *Client) authenticate(ctx context.Context) error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Check if token is still valid
	if c.token != "" && time.Now().Before(c.tokenExpiry.Add(-TokenRefreshBuffer)) {
		return nil
	}

	loginReq := LoginRequest{
		APIKey: c.apiKey,
		PIN:    c.pin,
	}

	var loginResp LoginResponse
	var errResp ErrorResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(loginReq).
		SetSuccessResult(&loginResp).
		SetErrorResult(&errResp).
		Post("/login")

	if err != nil {
		return fmt.Errorf("tvdb login request: %w", err)
	}

	if resp.IsErrorState() {
		return fmt.Errorf("tvdb login error %d: %s", resp.StatusCode, errResp.Message)
	}

	if loginResp.Data.Token == "" {
		return fmt.Errorf("tvdb login: empty token")
	}

	c.token = loginResp.Data.Token
	// TVDb tokens are valid for 30 days, but we refresh earlier
	c.tokenExpiry = time.Now().Add(24 * time.Hour)

	return nil
}

// request creates an authenticated request.
func (c *Client) request(ctx context.Context) (*req.Request, error) {
	// Ensure we have a valid token
	if err := c.authenticate(ctx); err != nil {
		return nil, err
	}

	c.tokenMutex.RLock()
	token := c.token
	c.tokenMutex.RUnlock()

	return c.httpClient.R().
		SetContext(ctx).
		SetBearerAuthToken(token), nil
}

// waitRateLimit waits for the rate limiter.
func (c *Client) waitRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

// getFromCache retrieves a value from cache.
func (c *Client) getFromCache(key string) any {
	val, ok := c.cache.Get(key)
	if !ok {
		return nil
	}
	return val
}

// setCache stores a value in cache.
func (c *Client) setCache(key string, data any, _ time.Duration) {
	c.cache.Set(key, data)
}

// clearCache removes all cached entries.
func (c *Client) clearCache() {
	c.cache.Clear()
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
func (c *Client) parseError(resp *req.Response, errResp *ErrorResponse) error {
	if errResp != nil && errResp.Message != "" {
		return fmt.Errorf("tvdb api error %d: %s", resp.StatusCode, errResp.Message)
	}
	return fmt.Errorf("tvdb api error: status %d", resp.StatusCode)
}

// Search performs a search query.
func (c *Client) Search(ctx context.Context, query string, searchType string) (*SearchResponse, error) {
	key := cacheKey("search", query, searchType)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SearchResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"query": query,
	}
	if searchType != "" {
		params["type"] = searchType
	}

	var result SearchResponse
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get("/search")

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result, SearchCacheTTL)
	return &result, nil
}

// GetSeries retrieves series details.
func (c *Client) GetSeries(ctx context.Context, id int) (*SeriesResponse, error) {
	key := cacheKey("series", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SeriesResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[SeriesResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/series/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetSeriesExtended retrieves series details with extended info.
func (c *Client) GetSeriesExtended(ctx context.Context, id int, meta string) (*SeriesExtendedResponse, error) {
	key := cacheKey("series:extended", id, meta)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SeriesExtendedResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if meta != "" {
		params["meta"] = meta
	}

	var result BaseResponse[SeriesExtendedResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/series/%d/extended", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetSeriesArtworks retrieves series artworks.
func (c *Client) GetSeriesArtworks(ctx context.Context, id int, artworkType *int, language string) ([]ArtworkResponse, error) {
	key := cacheKey("series:artworks", id, artworkType, language)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.([]ArtworkResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if artworkType != nil {
		params["type"] = fmt.Sprintf("%d", *artworkType)
	}
	if language != "" {
		params["lang"] = language
	}

	var result BaseResponse[struct {
		Artworks []ArtworkResponse `json:"artworks"`
	}]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/series/%d/artworks", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, result.Data.Artworks, 0)
	return result.Data.Artworks, nil
}

// GetSeriesEpisodes retrieves series episodes.
func (c *Client) GetSeriesEpisodes(ctx context.Context, id int, seasonType string, season *int, page int) ([]EpisodeResponse, error) {
	key := cacheKey("series:episodes", id, seasonType, season, page)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.([]EpisodeResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
	}
	if season != nil {
		params["season"] = fmt.Sprintf("%d", *season)
	}

	if seasonType == "" {
		seasonType = "default"
	}

	var result BaseResponse[struct {
		Episodes []EpisodeResponse `json:"episodes"`
	}]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/series/%d/episodes/%s", id, seasonType))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, result.Data.Episodes, 0)
	return result.Data.Episodes, nil
}

// GetSeason retrieves season details.
func (c *Client) GetSeason(ctx context.Context, id int) (*SeasonResponse, error) {
	key := cacheKey("season", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SeasonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[SeasonResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/seasons/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetSeasonExtended retrieves season details with extended info.
func (c *Client) GetSeasonExtended(ctx context.Context, id int) (*SeasonResponse, error) {
	key := cacheKey("season:extended", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*SeasonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[SeasonResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/seasons/%d/extended", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetEpisode retrieves episode details.
func (c *Client) GetEpisode(ctx context.Context, id int) (*EpisodeResponse, error) {
	key := cacheKey("episode", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*EpisodeResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[EpisodeResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/episodes/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetEpisodeExtended retrieves episode details with extended info.
func (c *Client) GetEpisodeExtended(ctx context.Context, id int, meta string) (*EpisodeResponse, error) {
	key := cacheKey("episode:extended", id, meta)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*EpisodeResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if meta != "" {
		params["meta"] = meta
	}

	var result BaseResponse[EpisodeResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/episodes/%d/extended", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetMovie retrieves movie details.
func (c *Client) GetMovie(ctx context.Context, id int) (*MovieResponse, error) {
	key := cacheKey("movie", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*MovieResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[MovieResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/movies/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetMovieExtended retrieves movie details with extended info.
func (c *Client) GetMovieExtended(ctx context.Context, id int, meta string) (*MovieResponse, error) {
	key := cacheKey("movie:extended", id, meta)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*MovieResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if meta != "" {
		params["meta"] = meta
	}

	var result BaseResponse[MovieResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/movies/%d/extended", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetPerson retrieves person details.
func (c *Client) GetPerson(ctx context.Context, id int) (*PersonResponse, error) {
	key := cacheKey("person", id)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	var result BaseResponse[PersonResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/people/%d", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// GetPersonExtended retrieves person details with extended info.
func (c *Client) GetPersonExtended(ctx context.Context, id int, meta string) (*PersonResponse, error) {
	key := cacheKey("person:extended", id, meta)
	if cached := c.getFromCache(key); cached != nil {
		if result, ok := cached.(*PersonResponse); ok {
			return result, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	req, err := c.request(ctx)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	if meta != "" {
		params["meta"] = meta
	}

	var result BaseResponse[PersonResponse]
	var errResp ErrorResponse

	resp, err := req.
		SetQueryParams(params).
		SetSuccessResult(&result).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("/people/%d/extended", id))

	if err != nil {
		return nil, fmt.Errorf("tvdb api request: %w", err)
	}

	if resp.IsErrorState() {
		return nil, c.parseError(resp, &errResp)
	}

	c.setCache(key, &result.Data, 0)
	return &result.Data, nil
}

// ClearCache clears all cached data.
func (c *Client) ClearCache() {
	c.clearCache()
}

// Logout invalidates the current token.
func (c *Client) Logout() {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.token = ""
	c.tokenExpiry = time.Time{}
}
