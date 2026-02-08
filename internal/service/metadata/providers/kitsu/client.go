package kitsu

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
	// BaseURL is the Kitsu API base URL.
	BaseURL = "https://kitsu.io/api/edge"

	// DefaultRateLimit is 2 requests per second (Kitsu has no published limit, be conservative).
	DefaultRateLimit = rate.Limit(2.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 5

	// DefaultCacheTTL is the cache duration.
	DefaultCacheTTL = 1 * time.Hour
)

// Config configures the Kitsu client.
type Config struct {
	// Enabled activates Kitsu as a metadata provider.
	Enabled bool

	// RateLimit is requests per second (default: 2).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 5).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is the Kitsu API client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new Kitsu client.
func NewClient(config Config) (*Client, error) {
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

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("create kitsu cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonHeader("Accept", "application/vnd.api+json").
		SetCommonHeader("Content-Type", "application/vnd.api+json")

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

// SearchAnime searches for anime by text query.
func (c *Client) SearchAnime(ctx context.Context, query string, limit, offset int) (*ListResponse[AnimeAttributes], error) {
	cacheKey := fmt.Sprintf("search:anime:%s:%d:%d", query, limit, offset)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*ListResponse[AnimeAttributes]); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result ListResponse[AnimeAttributes]
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("filter[text]", query).
		SetQueryParam("page[limit]", strconv.Itoa(limit)).
		SetQueryParam("page[offset]", strconv.Itoa(offset)).
		SetQueryParam("fields[anime]", "slug,synopsis,titles,canonicalTitle,averageRating,userCount,favoritesCount,startDate,endDate,popularityRank,ratingRank,ageRating,ageRatingGuide,subtype,status,posterImage,coverImage,episodeCount,episodeLength,youtubeVideoId,nsfw").
		SetSuccessResult(&result).
		Get("/anime")
	if err != nil {
		return nil, fmt.Errorf("kitsu search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("kitsu search: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetAnime fetches anime details by Kitsu ID.
func (c *Client) GetAnime(ctx context.Context, id string) (*SingleResponse[AnimeAttributes], error) {
	cacheKey := "anime:" + id
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*SingleResponse[AnimeAttributes]); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result SingleResponse[AnimeAttributes]
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("include", "categories,mappings").
		SetSuccessResult(&result).
		Get("/anime/" + id)
	if err != nil {
		return nil, fmt.Errorf("kitsu anime: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("kitsu anime: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetEpisodes fetches episodes for an anime.
func (c *Client) GetEpisodes(ctx context.Context, animeID string, limit, offset int) (*ListResponse[EpisodeAttributes], error) {
	cacheKey := fmt.Sprintf("episodes:%s:%d:%d", animeID, limit, offset)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*ListResponse[EpisodeAttributes]); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result ListResponse[EpisodeAttributes]
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("filter[mediaType]", "Anime").
		SetQueryParam("filter[media_id]", animeID).
		SetQueryParam("page[limit]", strconv.Itoa(limit)).
		SetQueryParam("page[offset]", strconv.Itoa(offset)).
		SetQueryParam("sort", "number").
		SetSuccessResult(&result).
		Get("/episodes")
	if err != nil {
		return nil, fmt.Errorf("kitsu episodes: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("kitsu episodes: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetMappings fetches external ID mappings for an anime.
func (c *Client) GetMappings(ctx context.Context, animeID string) (*ListResponse[MappingAttributes], error) {
	cacheKey := "mappings:" + animeID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*ListResponse[MappingAttributes]); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result ListResponse[MappingAttributes]
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("filter[externalSite]", "myanimelist/anime,thetvdb/series,thetvdb,anidb,anilist/anime").
		SetSuccessResult(&result).
		Get("/anime/" + animeID + "/mappings")
	if err != nil {
		return nil, fmt.Errorf("kitsu mappings: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("kitsu mappings: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetCategories fetches categories for an anime.
func (c *Client) GetCategories(ctx context.Context, animeID string) (*ListResponse[CategoryAttributes], error) {
	cacheKey := "categories:" + animeID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*ListResponse[CategoryAttributes]); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result ListResponse[CategoryAttributes]
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("page[limit]", "20").
		SetSuccessResult(&result).
		Get("/anime/" + animeID + "/categories")
	if err != nil {
		return nil, fmt.Errorf("kitsu categories: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("kitsu categories: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}
