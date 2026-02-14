package anilist

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
)

const (
	// BaseURL is the AniList GraphQL endpoint.
	BaseURL = "https://graphql.anilist.co"

	// DefaultRateLimit is 1.5 requests per second (AniList allows 90/min).
	DefaultRateLimit = rate.Limit(1.5)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 5

	// DefaultCacheTTL is the cache duration.
	DefaultCacheTTL = 1 * time.Hour
)

// Config configures the AniList client.
type Config struct {
	// Enabled activates AniList as a metadata provider.
	Enabled bool

	// RateLimit is requests per second (default: 1.5).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 5).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is the AniList GraphQL API client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new AniList client.
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

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		return nil, fmt.Errorf("create anilist cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonContentType("application/json").
		SetCommonHeader("Accept", "application/json").
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	// Circuit breaker
	circuitbreaker.WrapReqClient(client, "anilist", circuitbreaker.TierExternal)

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

// graphql executes a GraphQL query against the AniList API.
func (c *Client) graphql(ctx context.Context, query string, variables map[string]any, result any) error {
	if err := c.waitRateLimit(ctx); err != nil {
		return err
	}

	body := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	resp, err := c.httpClient.R().SetContext(ctx).
		SetBody(body).
		SetSuccessResult(result).
		Post("")
	if err != nil {
		return fmt.Errorf("anilist graphql: %w", err)
	}
	if resp.IsErrorState() {
		return fmt.Errorf("anilist graphql: status %d", resp.StatusCode)
	}

	return nil
}

// SearchAnime searches for anime by title.
func (c *Client) SearchAnime(ctx context.Context, query string, page, perPage int, isAdult bool) (*Page, error) {
	cacheKey := fmt.Sprintf("search:anime:%s:%d:%d:%t", query, page, perPage, isAdult)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Page); ok {
			return v, nil
		}
	}

	variables := map[string]any{
		"search":  query,
		"page":    page,
		"perPage": perPage,
		"type":    "ANIME",
	}
	if !isAdult {
		variables["isAdult"] = false
	}

	var resp GraphQLResponse[PageData]
	if err := c.graphql(ctx, searchQuery, variables, &resp); err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("anilist: %s", resp.Errors[0].Message)
	}

	c.setCache(cacheKey, &resp.Data.Page)
	return &resp.Data.Page, nil
}

// GetAnime fetches anime details by AniList ID.
func (c *Client) GetAnime(ctx context.Context, id int) (*Media, error) {
	cacheKey := fmt.Sprintf("anime:%d", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Media); ok {
			return v, nil
		}
	}

	variables := map[string]any{
		"id":   id,
		"type": "ANIME",
	}

	var resp GraphQLResponse[MediaData]
	if err := c.graphql(ctx, mediaQuery, variables, &resp); err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("anilist: %s", resp.Errors[0].Message)
	}
	if resp.Data.Media == nil {
		return nil, nil
	}

	c.setCache(cacheKey, resp.Data.Media)
	return resp.Data.Media, nil
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}

// GraphQL query strings.

const searchQuery = `
query ($search: String, $page: Int, $perPage: Int, $type: MediaType, $isAdult: Boolean) {
  Page(page: $page, perPage: $perPage) {
    pageInfo {
      total
      currentPage
      lastPage
      hasNextPage
      perPage
    }
    media(search: $search, type: $type, isAdult: $isAdult, sort: SEARCH_MATCH) {
      id
      idMal
      title {
        romaji
        english
        native
        userPreferred
      }
      type
      format
      status
      description(asHtml: false)
      startDate { year month day }
      endDate { year month day }
      season
      seasonYear
      episodes
      duration
      countryOfOrigin
      genres
      averageScore
      meanScore
      popularity
      favourites
      isAdult
      siteUrl
      coverImage {
        extraLarge
        large
        medium
        color
      }
      bannerImage
    }
  }
}
`

const mediaQuery = `
query ($id: Int, $type: MediaType) {
  Media(id: $id, type: $type) {
    id
    idMal
    title {
      romaji
      english
      native
      userPreferred
    }
    type
    format
    status
    description(asHtml: false)
    startDate { year month day }
    endDate { year month day }
    season
    seasonYear
    episodes
    duration
    countryOfOrigin
    source
    genres
    tags {
      id
      name
      category
      rank
      isAdult
    }
    averageScore
    meanScore
    popularity
    favourites
    isAdult
    siteUrl
    coverImage {
      extraLarge
      large
      medium
      color
    }
    bannerImage
    synonyms
    studios {
      edges {
        node { id name }
        isMain
      }
    }
    externalLinks {
      id
      url
      site
      siteId
      type
      language
    }
    trailer {
      id
      site
      thumbnail
    }
    characters(sort: [ROLE, RELEVANCE], perPage: 25) {
      edges {
        node {
          id
          name { full native }
          image { large medium }
          gender
        }
        role
        voiceActors(language: JAPANESE, sort: RELEVANCE) {
          id
          name { full native }
          languageV2
          image { large medium }
        }
      }
    }
    staff(sort: RELEVANCE, perPage: 25) {
      edges {
        node {
          id
          name { full native }
          image { large medium }
          primaryOccupations
        }
        role
      }
    }
    relations {
      edges {
        node {
          id
          title { romaji english native }
          type
          format
          status
          coverImage { large medium }
          averageScore
        }
        relationType
      }
    }
  }
}
`
