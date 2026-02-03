package metadata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

const (
	tmdbBaseURL  = "https://api.themoviedb.org/3"
	tmdbImageURL = "https://image.tmdb.org/t/p"
)

type TMDbClient struct {
	client      *resty.Client
	apiKey      string
	rateLimiter *rate.Limiter
	cache       sync.Map
	cacheTTL    time.Duration
}

type TMDbConfig struct {
	APIKey    string
	RateLimit rate.Limit
	CacheTTL  time.Duration
	ProxyURL  string
}

func NewTMDbClient(config TMDbConfig) *TMDbClient {
	client := resty.New().
		SetBaseURL(tmdbBaseURL).
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	if config.ProxyURL != "" {
		client.SetProxy(config.ProxyURL)
	}

	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(4.0)
	}

	if config.CacheTTL == 0 {
		config.CacheTTL = 24 * time.Hour
	}

	return &TMDbClient{
		client:      client,
		apiKey:      config.APIKey,
		rateLimiter: rate.NewLimiter(config.RateLimit, 10),
		cacheTTL:    config.CacheTTL,
	}
}

func (c *TMDbClient) SearchMovies(ctx context.Context, query string, year *int) (*TMDbSearchResponse, error) {
	cacheKey := fmt.Sprintf("search:%s:%v", query, year)
	if cached := c.getFromCache(cacheKey); cached != nil {
		return cached.(*TMDbSearchResponse), nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	params := map[string]string{
		"api_key": c.apiKey,
		"query":   query,
	}

	if year != nil {
		params["year"] = fmt.Sprintf("%d", *year)
	}

	var result TMDbSearchResponse
	var errResp TMDbError

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&result).
		SetError(&errResp).
		Get("/search/movie")

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp.StatusCode(), &errResp)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

func (c *TMDbClient) GetMovie(ctx context.Context, tmdbID int) (*TMDbMovie, error) {
	cacheKey := fmt.Sprintf("movie:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		return cached.(*TMDbMovie), nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result TMDbMovie
	var errResp TMDbError

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("api_key", c.apiKey).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d", tmdbID))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp.StatusCode(), &errResp)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

func (c *TMDbClient) GetMovieCredits(ctx context.Context, tmdbID int) (*TMDbCredits, error) {
	cacheKey := fmt.Sprintf("credits:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		return cached.(*TMDbCredits), nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result TMDbCredits
	var errResp TMDbError

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("api_key", c.apiKey).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/credits", tmdbID))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp.StatusCode(), &errResp)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

func (c *TMDbClient) GetMovieImages(ctx context.Context, tmdbID int) (*TMDbImages, error) {
	cacheKey := fmt.Sprintf("images:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		return cached.(*TMDbImages), nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result TMDbImages
	var errResp TMDbError

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("api_key", c.apiKey).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/movie/%d/images", tmdbID))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp.StatusCode(), &errResp)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

func (c *TMDbClient) GetCollection(ctx context.Context, collectionID int) (*TMDbCollectionDetails, error) {
	cacheKey := fmt.Sprintf("collection:%d", collectionID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		return cached.(*TMDbCollectionDetails), nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result TMDbCollectionDetails
	var errResp TMDbError

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("api_key", c.apiKey).
		SetResult(&result).
		SetError(&errResp).
		Get(fmt.Sprintf("/collection/%d", collectionID))

	if err != nil {
		return nil, fmt.Errorf("tmdb api request: %w", err)
	}

	if resp.IsError() {
		return nil, c.parseError(resp.StatusCode(), &errResp)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

func (c *TMDbClient) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s%s", tmdbImageURL, size, path)
}

func (c *TMDbClient) DownloadImage(ctx context.Context, path string, size string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("empty image path")
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	url := c.GetImageURL(path, size)

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

func (c *TMDbClient) getFromCache(key string) interface{} {
	if val, ok := c.cache.Load(key); ok {
		entry := val.(*CacheEntry)
		if !entry.IsExpired() {
			return entry.Data
		}
		c.cache.Delete(key)
	}
	return nil
}

func (c *TMDbClient) setCache(key string, data interface{}) {
	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.cacheTTL),
	}
	c.cache.Store(key, entry)
}

func (c *TMDbClient) parseError(statusCode int, errResp *TMDbError) error {
	if errResp != nil && errResp.StatusMessage != "" {
		return fmt.Errorf("tmdb api error %d: %s", statusCode, errResp.StatusMessage)
	}
	return fmt.Errorf("tmdb api error: status %d", statusCode)
}

func (c *TMDbClient) ClearCache() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})
}
