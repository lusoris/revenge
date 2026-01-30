package whisparr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lusoris/revenge/pkg/resilience"
)

// Client errors.
var (
	ErrNotFound      = errors.New("movie not found")
	ErrUnauthorized  = errors.New("invalid API key")
	ErrUnavailable   = errors.New("whisparr unavailable")
	ErrNotConfigured = errors.New("whisparr not configured")
)

// ClientConfig holds Whisparr client configuration.
type ClientConfig struct {
	BaseURL    string        `koanf:"base_url"`
	APIKey     string        `koanf:"api_key"`
	Timeout    time.Duration `koanf:"timeout"`
	RetryCount int           `koanf:"retry_count"`
}

// DefaultClientConfig returns sensible defaults.
var DefaultClientConfig = ClientConfig{
	Timeout:    30 * time.Second,
	RetryCount: 3,
}

// Client is a Whisparr API v3 client with resilience patterns.
// Whisparr is a fork of Radarr for adult content management.
type Client struct {
	http    *resty.Client
	breaker *resilience.CircuitBreaker
	config  ClientConfig
	logger  *slog.Logger
}

// NewClient creates a new Whisparr client.
func NewClient(cfg ClientConfig, logger *slog.Logger) (*Client, error) {
	if cfg.BaseURL == "" || cfg.APIKey == "" {
		return nil, ErrNotConfigured
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultClientConfig.Timeout
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = DefaultClientConfig.RetryCount
	}

	// Create resty client
	http := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Api-Key", cfg.APIKey).
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true
			}
			return r.StatusCode() >= 500
		}).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			logger.Debug("whisparr request", "method", r.Method, "url", r.URL)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Debug("whisparr response", "status", r.StatusCode())
			return nil
		})

	// Create circuit breaker
	breaker := resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
		Name:                "whisparr",
		MaxFailures:         5,
		Timeout:             30 * time.Second,
		MaxHalfOpenRequests: 3,
		OnStateChange: func(name string, from, to int) {
			states := []string{"closed", "open", "half-open"}
			logger.Warn("circuit breaker state change",
				"name", name,
				"from", states[from],
				"to", states[to],
			)
		},
		IsSuccessful: func(err error) bool {
			return err == nil || errors.Is(err, ErrNotFound)
		},
	})

	return &Client{
		http:    http,
		breaker: breaker,
		config:  cfg,
		logger:  logger.With("client", "whisparr"),
	}, nil
}

// GetSystemStatus retrieves Whisparr system status.
func (c *Client) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	var status SystemStatus

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&status).
			Get("/api/v3/system/status")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetHealth retrieves health check results.
func (c *Client) GetHealth(ctx context.Context) ([]HealthCheck, error) {
	var checks []HealthCheck

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&checks).
			Get("/api/v3/health")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return checks, nil
}

// ListMovies retrieves all movies/scenes from Whisparr.
func (c *Client) ListMovies(ctx context.Context) ([]Movie, error) {
	var movies []Movie

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&movies).
			Get("/api/v3/movie")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return movies, nil
}

// GetMovie retrieves a movie by Whisparr ID.
func (c *Client) GetMovie(ctx context.Context, id int) (*Movie, error) {
	var movie Movie

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&movie).
			SetPathParam("id", fmt.Sprint(id)).
			Get("/api/v3/movie/{id}")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetMovieByForeignID retrieves a movie by external ID (StashDB/TPDB).
func (c *Client) GetMovieByForeignID(ctx context.Context, foreignID string) (*Movie, error) {
	var movies []Movie

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&movies).
			SetQueryParam("foreignId", foreignID).
			Get("/api/v3/movie")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	if len(movies) == 0 {
		return nil, ErrNotFound
	}

	return &movies[0], nil
}

// GetMovieByStashID retrieves a movie by StashDB scene ID.
func (c *Client) GetMovieByStashID(ctx context.Context, stashID string) (*Movie, error) {
	movies, err := c.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range movies {
		if m.StashID == stashID || m.ForeignId == stashID {
			return &m, nil
		}
	}

	return nil, ErrNotFound
}

// SearchMovies searches for movies by query.
func (c *Client) SearchMovies(ctx context.Context, query string) ([]SearchMovieResult, error) {
	var results []SearchMovieResult

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&results).
			SetQueryParam("term", query).
			Get("/api/v3/movie/lookup")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// LookupMovieByForeignID looks up a movie by external ID.
func (c *Client) LookupMovieByForeignID(ctx context.Context, foreignID string) (*SearchMovieResult, error) {
	var result SearchMovieResult

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&result).
			SetQueryParam("foreignId", foreignID).
			Get("/api/v3/movie/lookup/id")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// AddMovie adds a movie to Whisparr.
func (c *Client) AddMovie(ctx context.Context, opts AddMovieOptions) (*Movie, error) {
	var movie Movie

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(opts).
			SetResult(&movie).
			Post("/api/v3/movie")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// UpdateMovie updates a movie in Whisparr.
func (c *Client) UpdateMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	var updated Movie

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(movie).
			SetResult(&updated).
			SetPathParam("id", fmt.Sprint(movie.ID)).
			Put("/api/v3/movie/{id}")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteMovie removes a movie from Whisparr.
func (c *Client) DeleteMovie(ctx context.Context, id int, deleteFiles bool) error {
	return c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetPathParam("id", fmt.Sprint(id)).
			SetQueryParam("deleteFiles", fmt.Sprint(deleteFiles)).
			Delete("/api/v3/movie/{id}")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})
}

// ListPerformers retrieves all performers.
func (c *Client) ListPerformers(ctx context.Context) ([]Performer, error) {
	var performers []Performer

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&performers).
			Get("/api/v3/person")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return performers, nil
}

// GetPerformer retrieves a performer by ID.
func (c *Client) GetPerformer(ctx context.Context, id int) (*Performer, error) {
	var performer Performer

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&performer).
			SetPathParam("id", fmt.Sprint(id)).
			Get("/api/v3/person/{id}")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return &performer, nil
}

// ListRootFolders retrieves root folder configurations.
func (c *Client) ListRootFolders(ctx context.Context) ([]RootFolder, error) {
	var folders []RootFolder

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&folders).
			Get("/api/v3/rootfolder")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return folders, nil
}

// ListQualityProfiles retrieves quality profiles.
func (c *Client) ListQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	var profiles []QualityProfile

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&profiles).
			Get("/api/v3/qualityprofile")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// RescanMovie triggers a rescan of a movie's files.
func (c *Client) RescanMovie(ctx context.Context, movieID int) error {
	return c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		body := map[string]any{
			"name":    "RescanMovie",
			"movieId": movieID,
		}

		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(body).
			Post("/api/v3/command")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})
}

// RefreshMovie triggers a metadata refresh for a movie.
func (c *Client) RefreshMovie(ctx context.Context, movieID int) error {
	return c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		body := map[string]any{
			"name":    "RefreshMovie",
			"movieId": movieID,
		}

		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(body).
			Post("/api/v3/command")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})
}

// SearchMovieByID triggers a search for a movie.
func (c *Client) SearchMovieByID(ctx context.Context, movieID int) error {
	return c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		body := map[string]any{
			"name":     "MoviesSearch",
			"movieIds": []int{movieID},
		}

		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(body).
			Post("/api/v3/command")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})
}

// ImportMovie imports a movie from a path.
func (c *Client) ImportMovie(ctx context.Context, path string, movieID int) error {
	return c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		body := map[string]any{
			"name": "ManualImport",
			"files": []map[string]any{
				{
					"path":    path,
					"movieId": movieID,
				},
			},
			"importMode": "auto",
		}

		resp, err := c.http.R().
			SetContext(ctx).
			SetBody(body).
			Post("/api/v3/command")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp)
	})
}

// handleResponse processes HTTP responses and converts to errors.
func (c *Client) handleResponse(resp *resty.Response) error {
	switch resp.StatusCode() {
	case 200, 201, 202:
		return nil
	case 401:
		return ErrUnauthorized
	case 404:
		return ErrNotFound
	case 500, 502, 503, 504:
		return ErrUnavailable
	default:
		return fmt.Errorf("unexpected status: %d - %s", resp.StatusCode(), resp.String())
	}
}

// Ping checks if Whisparr is reachable.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.GetSystemStatus(ctx)
	return err
}

// IsConfigured returns true if the client is configured.
func (c *Client) IsConfigured() bool {
	return c.config.BaseURL != "" && c.config.APIKey != ""
}

// Stats returns circuit breaker statistics.
func (c *Client) Stats() ClientStats {
	cbStats := c.breaker.Stats()
	states := []string{"closed", "open", "half-open"}

	return ClientStats{
		CircuitBreakerState:    states[cbStats.State],
		CircuitBreakerFailures: cbStats.Failures,
		CircuitBreakerRequests: cbStats.Requests,
	}
}

// ClientStats contains client statistics.
type ClientStats struct {
	CircuitBreakerState    string
	CircuitBreakerFailures int
	CircuitBreakerRequests int
}

// Close releases client resources.
func (c *Client) Close() error {
	return nil
}
