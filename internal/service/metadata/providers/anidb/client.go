package anidb

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
)

const (
	// BaseURL is the AniDB HTTP API endpoint.
	BaseURL = "http://api.anidb.net:9001/httpapi"

	// ImageBaseURL is the AniDB image CDN.
	ImageBaseURL = "https://cdn.anidb.net/images/main/"

	// TitleDumpURL is the daily anime title dump.
	TitleDumpURL = "http://anidb.net/api/anime-titles.dat.gz"

	// DefaultRateLimit is 1 request every 2 seconds (AniDB limit).
	DefaultRateLimit = rate.Limit(0.5)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 1

	// DefaultCacheTTL is the cache duration.
	DefaultCacheTTL = 2 * time.Hour

	// ProtocolVersion is the AniDB HTTP API protocol version.
	ProtocolVersion = "1"
)

// Config configures the AniDB client.
type Config struct {
	// Enabled activates AniDB as a metadata provider.
	Enabled bool

	// ClientName is the registered AniDB client identifier (required).
	// Register at https://anidb.net/perl-bin/animedb.pl?show=client
	ClientName string

	// ClientVersion is the client version number (default: 1).
	ClientVersion int

	// RateLimit is requests per second (default: 0.5 = 1 req/2s).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 1).
	Burst int

	// CacheTTL is the cache duration (default: 2h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 30s).
	Timeout time.Duration
}

// Client is the AniDB HTTP API client with rate limiting and caching.
type Client struct {
	httpClient    *req.Client
	rawClient     *http.Client
	rateLimiter   *rate.Limiter
	cache         *cache.L1Cache[string, any]
	cacheTTL      time.Duration
	clientName    string
	clientVersion int

	// Title index for search
	titleIndex    []TitleDumpEntry
	titleIndexMu  sync.RWMutex
	lastTitleLoad time.Time
}

// NewClient creates a new AniDB client.
func NewClient(config Config) (*Client, error) {
	if config.ClientName == "" {
		return nil, fmt.Errorf("anidb: client_name is required (register at https://anidb.net/perl-bin/animedb.pl?show=client)")
	}
	if config.ClientVersion == 0 {
		config.ClientVersion = 1
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
		config.Timeout = 30 * time.Second
	}

	l1, err := cache.NewL1Cache[string, any](5000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		return nil, fmt.Errorf("create anidb cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(1).
		SetCommonRetryBackoffInterval(3*time.Second, 10*time.Second).
		SetCommonHeader("Accept", "application/xml").
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	// Circuit breaker
	circuitbreaker.WrapReqClient(client, "anidb", circuitbreaker.TierExternal)

	return &Client{
		httpClient:    client,
		rawClient:     &http.Client{Timeout: config.Timeout},
		rateLimiter:   rate.NewLimiter(config.RateLimit, config.Burst),
		cache:         l1,
		cacheTTL:      config.CacheTTL,
		clientName:    config.ClientName,
		clientVersion: config.ClientVersion,
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

// GetAnime fetches anime details by AniDB ID.
func (c *Client) GetAnime(ctx context.Context, aid int) (*AnimeResponse, error) {
	cacheKey := fmt.Sprintf("anime:%d", aid)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*AnimeResponse); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("request", "anime").
		SetQueryParam("aid", strconv.Itoa(aid)).
		SetQueryParam("client", c.clientName).
		SetQueryParam("clientver", strconv.Itoa(c.clientVersion)).
		SetQueryParam("protover", ProtocolVersion).
		Get("")
	if err != nil {
		return nil, fmt.Errorf("anidb anime: %w", err)
	}

	body := resp.Bytes()

	// Check for error response
	var errResp ErrorResponse
	if xml.Unmarshal(body, &errResp) == nil && errResp.Text != "" {
		return nil, fmt.Errorf("anidb: %s", errResp.Text)
	}

	var anime AnimeResponse
	if err := xml.Unmarshal(body, &anime); err != nil {
		return nil, fmt.Errorf("anidb: parse xml: %w", err)
	}

	c.setCache(cacheKey, &anime)
	return &anime, nil
}

// SearchAnime searches the title dump index for matching anime.
func (c *Client) SearchAnime(ctx context.Context, query string, limit int) ([]TitleDumpEntry, error) {
	if err := c.ensureTitleIndex(ctx); err != nil {
		return nil, fmt.Errorf("anidb search: %w", err)
	}

	c.titleIndexMu.RLock()
	defer c.titleIndexMu.RUnlock()

	queryLower := strings.ToLower(query)
	var matches []TitleDumpEntry
	seen := make(map[int]bool)

	// Exact matches first
	for _, entry := range c.titleIndex {
		if strings.EqualFold(entry.Title, query) && !seen[entry.AID] {
			matches = append(matches, entry)
			seen[entry.AID] = true
			if len(matches) >= limit {
				return matches, nil
			}
		}
	}

	// Prefix matches
	for _, entry := range c.titleIndex {
		if strings.HasPrefix(strings.ToLower(entry.Title), queryLower) && !seen[entry.AID] {
			matches = append(matches, entry)
			seen[entry.AID] = true
			if len(matches) >= limit {
				return matches, nil
			}
		}
	}

	// Contains matches
	for _, entry := range c.titleIndex {
		if strings.Contains(strings.ToLower(entry.Title), queryLower) && !seen[entry.AID] {
			matches = append(matches, entry)
			seen[entry.AID] = true
			if len(matches) >= limit {
				return matches, nil
			}
		}
	}

	return matches, nil
}

// ensureTitleIndex loads the title dump if not loaded or stale (>24h).
func (c *Client) ensureTitleIndex(ctx context.Context) error {
	c.titleIndexMu.RLock()
	if len(c.titleIndex) > 0 && time.Since(c.lastTitleLoad) < 24*time.Hour {
		c.titleIndexMu.RUnlock()
		return nil
	}
	c.titleIndexMu.RUnlock()

	return c.loadTitleDump(ctx)
}

// loadTitleDump downloads and parses the AniDB anime titles dump.
func (c *Client) loadTitleDump(ctx context.Context) error {
	c.titleIndexMu.Lock()
	defer c.titleIndexMu.Unlock()

	// Double-check after acquiring write lock
	if len(c.titleIndex) > 0 && time.Since(c.lastTitleLoad) < 24*time.Hour {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, TitleDumpURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.rawClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch title dump: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // best-effort cleanup

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("title dump: status %d", resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("title dump gzip: %w", err)
	}
	defer gz.Close() //nolint:errcheck // best-effort cleanup

	entries, err := parseTitleDump(gz)
	if err != nil {
		return err
	}

	c.titleIndex = entries
	c.lastTitleLoad = time.Now()
	return nil
}

// parseTitleDump parses the pipe-delimited title dump format.
// Format: aid|type|lang|title
func parseTitleDump(r io.Reader) ([]TitleDumpEntry, error) {
	var entries []TitleDumpEntry
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		aid, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		entries = append(entries, TitleDumpEntry{
			AID:   aid,
			Type:  parts[1],
			Lang:  parts[2],
			Title: parts[3],
		})
	}

	return entries, scanner.Err()
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}
