package letterboxd

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	baseURL  = "https://api.letterboxd.com/api/v0"
	tokenURL = baseURL + "/auth/token"
)

// Config contains configuration for the Letterboxd client.
type Config struct {
	Enabled   bool
	APIKey    string
	APISecret string
	RateLimit rate.Limit
	Burst     int
	CacheTTL  time.Duration
	Timeout   time.Duration
}

// Client is a client for the Letterboxd API.
type Client struct {
	client      *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	apiKey      string
	apiSecret   string
	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

// NewClient creates a new Letterboxd API client.
func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("letterboxd API key and secret are required")
	}
	if config.RateLimit <= 0 {
		config.RateLimit = 1.0
	}
	if config.Burst <= 0 {
		config.Burst = 3
	}
	if config.CacheTTL <= 0 {
		config.CacheTTL = time.Hour
	}
	if config.Timeout <= 0 {
		config.Timeout = 15 * time.Second
	}

	httpClient := req.C().
		SetBaseURL(baseURL).
		SetTimeout(config.Timeout).
		SetCommonHeader("Accept", "application/json").
		SetCommonRetryCount(2).
		SetCommonRetryFixedInterval(time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	l1, err := cache.NewL1Cache[string, any](2000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		l1, _ = cache.NewL1Cache[string, any](0, 0)
	}

	c := &Client{
		client:      httpClient,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cache:       l1,
		apiKey:      config.APIKey,
		apiSecret:   config.APISecret,
	}

	return c, nil
}

// ensureToken obtains or refreshes the OAuth2 access token.
func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-30*time.Second)) {
		return nil
	}

	var tokenResp TokenResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     c.apiKey,
			"client_secret": c.apiSecret,
		}).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetSuccessResult(&tokenResp).
		Post(tokenURL)

	if err != nil {
		return fmt.Errorf("letterboxd token request: %w", err)
	}
	if resp.IsErrorState() {
		return fmt.Errorf("letterboxd token request failed: %s", resp.Status)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	c.client.SetCommonBearerAuthToken(c.accessToken)

	return nil
}

// doGet performs an authenticated, rate-limited GET request.
func (c *Client) doGet(ctx context.Context, path string, queryParams map[string]string, result any) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	r := c.client.R().SetContext(ctx)
	if result != nil {
		r.SetSuccessResult(result)
	}
	for k, v := range queryParams {
		r.SetQueryParam(k, v)
	}

	resp, err := r.Get(path)
	if err != nil {
		return fmt.Errorf("letterboxd GET %s: %w", path, err)
	}
	if resp.IsErrorState() {
		return fmt.Errorf("letterboxd GET %s: %s", path, resp.Status)
	}

	return nil
}

func (c *Client) clearCache() {
	c.cache.Clear()
}

// SearchFilms searches for films by query string.
func (c *Client) SearchFilms(ctx context.Context, query string) (*SearchResponse, error) {
	cacheKey := "search:" + query
	if v, ok := c.cache.Get(cacheKey); ok {
		if r, ok := v.(*SearchResponse); ok {
			return r, nil
		}
	}

	var result SearchResponse
	err := c.doGet(ctx, "/search", map[string]string{
		"input":   query,
		"include": "FilmSearchItem",
		"perPage": "20",
	}, &result)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, &result)
	return &result, nil
}

// GetFilm gets detailed information about a film by Letterboxd ID.
// Supports LID, "tmdb:{id}", or "imdb:{id}" prefixes.
func (c *Client) GetFilm(ctx context.Context, id string) (*Film, error) {
	cacheKey := "film:" + id
	if v, ok := c.cache.Get(cacheKey); ok {
		if f, ok := v.(*Film); ok {
			return f, nil
		}
	}

	var film Film
	err := c.doGet(ctx, "/film/"+url.PathEscape(id), nil, &film)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, &film)
	return &film, nil
}

// GetFilmStatistics gets statistical data about a film.
func (c *Client) GetFilmStatistics(ctx context.Context, id string) (*FilmStatistics, error) {
	cacheKey := "film-stats:" + id
	if v, ok := c.cache.Get(cacheKey); ok {
		if s, ok := v.(*FilmStatistics); ok {
			return s, nil
		}
	}

	var stats FilmStatistics
	err := c.doGet(ctx, "/film/"+url.PathEscape(id)+"/statistics", nil, &stats)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, &stats)
	return &stats, nil
}

// GetContributor gets information about a film contributor.
// Supports LID or "tmdb:{id}" prefix.
func (c *Client) GetContributor(ctx context.Context, id string) (*Contributor, error) {
	cacheKey := "contributor:" + id
	if v, ok := c.cache.Get(cacheKey); ok {
		if ct, ok := v.(*Contributor); ok {
			return ct, nil
		}
	}

	var contributor Contributor
	err := c.doGet(ctx, "/contributor/"+url.PathEscape(id), nil, &contributor)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, &contributor)
	return &contributor, nil
}

// GetFilmCollection gets information about a film collection.
func (c *Client) GetFilmCollection(ctx context.Context, id string) (*FilmCollection, error) {
	cacheKey := "collection:" + id
	if v, ok := c.cache.Get(cacheKey); ok {
		if fc, ok := v.(*FilmCollection); ok {
			return fc, nil
		}
	}

	var collection FilmCollection
	err := c.doGet(ctx, "/film-collection/"+url.PathEscape(id), nil, &collection)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, &collection)
	return &collection, nil
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}
