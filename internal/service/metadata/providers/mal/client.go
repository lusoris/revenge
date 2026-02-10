package mal

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	// BaseURL is the MyAnimeList API v2 base URL.
	BaseURL = "https://api.myanimelist.net/v2"

	// DefaultRateLimit is 1 request per second (conservative; MAL rate limits undocumented).
	DefaultRateLimit = rate.Limit(1.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 3

	// DefaultCacheTTL is the cache duration.
	DefaultCacheTTL = 1 * time.Hour

	// searchFields are the fields requested in search results.
	searchFields = "id,title,main_picture,alternative_titles,start_date,end_date,synopsis,mean,popularity,num_list_users,media_type,status,genres,num_episodes,start_season,nsfw"

	// detailFields are the fields requested for full anime details.
	detailFields = "id,title,main_picture,alternative_titles,start_date,end_date,synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,media_type,status,genres,num_episodes,start_season,broadcast,source,average_episode_duration,rating,pictures,background,related_anime,recommendations,studios,statistics"
)

// Config configures the MAL client.
type Config struct {
	// Enabled activates MAL as a metadata provider.
	Enabled bool

	// ClientID is the MyAnimeList API client ID (required).
	// Obtain from https://myanimelist.net/apiconfig
	ClientID string

	// RateLimit is requests per second (default: 1.0).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 3).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is the MyAnimeList API v2 client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new MAL client.
func NewClient(config Config) (*Client, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf("mal: client_id is required")
	}
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
		config.Timeout = 15 * time.Second
	}

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		return nil, fmt.Errorf("create mal cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonHeader("Accept", "application/json").
		SetCommonHeader("X-MAL-CLIENT-ID", config.ClientID).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	return &Client{
		httpClient:  client,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cache:       l1,
		cacheTTL:    config.CacheTTL,
	}, nil
}

func (c *Client) waitRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

func (c *Client) getFromCache(key string) any {
	val, ok := c.cache.Get(key)
	if !ok {
		return nil
	}
	return val
}

func (c *Client) setCache(key string, data any) {
	c.cache.Set(key, data)
}

func (c *Client) clearCache() {
	c.cache.Clear()
}

// SearchAnime searches for anime by query string.
func (c *Client) SearchAnime(ctx context.Context, query string, limit, offset int) (*ListResponse, error) {
	cacheKey := fmt.Sprintf("search:anime:%s:%d:%d", query, limit, offset)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*ListResponse); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result ListResponse
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("limit", strconv.Itoa(limit)).
		SetQueryParam("offset", strconv.Itoa(offset)).
		SetQueryParam("fields", searchFields).
		SetSuccessResult(&result).
		Get("/anime")
	if err != nil {
		return nil, fmt.Errorf("mal search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("mal search: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetAnime fetches anime details by MAL ID.
func (c *Client) GetAnime(ctx context.Context, id int) (*Anime, error) {
	cacheKey := fmt.Sprintf("anime:%d", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Anime); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result Anime
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("fields", detailFields).
		SetSuccessResult(&result).
		Get("/anime/" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("mal anime: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("mal anime: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}
