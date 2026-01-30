package stashdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default StashDB API endpoint.
	DefaultBaseURL = "https://stashdb.org/graphql"

	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 30 * time.Second
)

// Errors returned by the client.
var (
	ErrNotFound     = errors.New("scene not found")
	ErrUnavailable  = errors.New("stashdb unavailable")
	ErrUnauthorized = errors.New("stashdb authentication required")
	ErrRateLimited  = errors.New("stashdb rate limited")
)

// ClientConfig contains configuration for the StashDB client.
type ClientConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// Client is a StashDB GraphQL API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new StashDB client.
func NewClient(cfg ClientConfig) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	return &Client{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// IsConfigured returns true if the client has an API key.
func (c *Client) IsConfigured() bool {
	return c.apiKey != ""
}

// Ping checks if the StashDB API is reachable.
func (c *Client) Ping(ctx context.Context) error {
	query := `query { version { version } }`
	var result struct {
		Version struct {
			Version string `json:"version"`
		} `json:"version"`
	}

	if err := c.query(ctx, query, nil, &result); err != nil {
		return err
	}
	return nil
}

// GetScene retrieves a scene by ID.
func (c *Client) GetScene(ctx context.Context, id string) (*Scene, error) {
	query := `
		query FindScene($id: ID!) {
			findScene(id: $id) {
				id
				title
				details
				date
				release_date
				production_date
				duration
				director
				code
				studio {
					id
					name
					urls { url type }
					images { id url width height }
				}
				performers {
					performer {
						id
						name
						disambiguation
						aliases
						gender
						birth_date
						country
						images { id url width height }
					}
					as
				}
				tags { id name description }
				images { id url width height }
				fingerprints { algorithm hash duration }
				urls { url type }
				created
				updated
			}
		}
	`

	var result struct {
		FindScene *Scene `json:"findScene"`
	}

	if err := c.query(ctx, query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}

	if result.FindScene == nil {
		return nil, ErrNotFound
	}

	return result.FindScene, nil
}

// SearchScenes searches for scenes by title.
func (c *Client) SearchScenes(ctx context.Context, query string, page, perPage int) (*SceneQueryResult, error) {
	gql := `
		query SearchScenes($query: String!, $page: Int!, $per_page: Int!) {
			searchScene(term: $query, page: $page, per_page: $per_page) {
				count
				scenes {
					id
					title
					details
					date
					duration
					studio { id name }
					performers {
						performer { id name }
						as
					}
					images { id url width height }
				}
			}
		}
	`

	var result struct {
		SearchScene struct {
			Count  int     `json:"count"`
			Scenes []Scene `json:"scenes"`
		} `json:"searchScene"`
	}

	vars := map[string]any{
		"query":    query,
		"page":     page,
		"per_page": perPage,
	}

	if err := c.query(ctx, gql, vars, &result); err != nil {
		return nil, err
	}

	return &SceneQueryResult{
		Count: result.SearchScene.Count,
		Data:  result.SearchScene.Scenes,
	}, nil
}

// FindSceneByFingerprint finds a scene by fingerprint (phash, oshash, or md5).
func (c *Client) FindSceneByFingerprint(ctx context.Context, algorithm, hash string, duration int) ([]Scene, error) {
	gql := `
		query FindSceneByFingerprint($fingerprint: FingerprintQueryInput!) {
			findSceneByFingerprint(fingerprint: $fingerprint) {
				id
				title
				details
				date
				duration
				director
				studio { id name }
				performers {
					performer { id name }
					as
				}
				images { id url width height }
				fingerprints { algorithm hash duration }
			}
		}
	`

	var result struct {
		FindSceneByFingerprint []Scene `json:"findSceneByFingerprint"`
	}

	vars := map[string]any{
		"fingerprint": map[string]any{
			"algorithm": algorithm,
			"hash":      hash,
			"duration":  duration,
		},
	}

	if err := c.query(ctx, gql, vars, &result); err != nil {
		return nil, err
	}

	return result.FindSceneByFingerprint, nil
}

// FindScenesByFingerprints finds scenes matching any of the given fingerprints.
func (c *Client) FindScenesByFingerprints(ctx context.Context, fingerprints []Fingerprint) ([]Scene, error) {
	gql := `
		query FindScenesByFingerprints($fingerprints: [FingerprintQueryInput!]!) {
			findScenesByFingerprints(fingerprints: $fingerprints) {
				id
				title
				details
				date
				duration
				director
				studio { id name }
				performers {
					performer { id name }
					as
				}
				images { id url width height }
				fingerprints { algorithm hash duration }
			}
		}
	`

	var result struct {
		FindScenesByFingerprints []Scene `json:"findScenesByFingerprints"`
	}

	fps := make([]map[string]any, len(fingerprints))
	for i, fp := range fingerprints {
		fps[i] = map[string]any{
			"algorithm": fp.Algorithm,
			"hash":      fp.Hash,
			"duration":  fp.Duration,
		}
	}

	vars := map[string]any{"fingerprints": fps}

	if err := c.query(ctx, gql, vars, &result); err != nil {
		return nil, err
	}

	return result.FindScenesByFingerprints, nil
}

// GetPerformer retrieves a performer by ID.
func (c *Client) GetPerformer(ctx context.Context, id string) (*Performer, error) {
	query := `
		query FindPerformer($id: ID!) {
			findPerformer(id: $id) {
				id
				name
				disambiguation
				aliases
				gender
				birth_date
				death_date
				age
				ethnicity
				country
				eye_color
				hair_color
				height
				cup_size
				band_size
				waist_size
				hip_size
				breast_type
				career_start_year
				career_end_year
				images { id url width height }
				scene_count
			}
		}
	`

	var result struct {
		FindPerformer *Performer `json:"findPerformer"`
	}

	if err := c.query(ctx, query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}

	if result.FindPerformer == nil {
		return nil, ErrNotFound
	}

	return result.FindPerformer, nil
}

// SearchPerformers searches for performers by name.
func (c *Client) SearchPerformers(ctx context.Context, query string, page, perPage int) (*PerformerQueryResult, error) {
	gql := `
		query SearchPerformers($query: String!, $page: Int!, $per_page: Int!) {
			searchPerformer(term: $query, page: $page, per_page: $per_page) {
				count
				performers {
					id
					name
					disambiguation
					aliases
					gender
					birth_date
					country
					images { id url width height }
					scene_count
				}
			}
		}
	`

	var result struct {
		SearchPerformer struct {
			Count      int         `json:"count"`
			Performers []Performer `json:"performers"`
		} `json:"searchPerformer"`
	}

	vars := map[string]any{
		"query":    query,
		"page":     page,
		"per_page": perPage,
	}

	if err := c.query(ctx, gql, vars, &result); err != nil {
		return nil, err
	}

	return &PerformerQueryResult{
		Count: result.SearchPerformer.Count,
		Data:  result.SearchPerformer.Performers,
	}, nil
}

// query executes a GraphQL query.
func (c *Client) query(ctx context.Context, query string, variables map[string]any, result any) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.apiKey != "" {
		req.Header.Set("ApiKey", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusOK:
		// Continue processing
	default:
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var gqlResp GraphQLResponse[json.RawMessage]
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return gqlResp.Errors[0]
	}

	if err := json.Unmarshal(gqlResp.Data, result); err != nil {
		return fmt.Errorf("decode data: %w", err)
	}

	return nil
}
