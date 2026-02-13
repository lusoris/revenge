package anidb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// AniDB is a comprehensive anime database with detailed episode data,
// character/seiyuu info, tags, and cross-references to MAL/ANN.
// It requires a registered client identifier (no API key, no OAuth).
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for AniDB.
type Provider struct {
	metadata.TVShowProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new AniDB provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create anidb provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 44, // Below AniList (45), above MAL (43)/Kitsu (42)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID { return metadata.ProviderAniDB }
func (p *Provider) Name() string            { return "AniDB" }
func (p *Provider) Priority() int           { return p.priority }
func (p *Provider) SupportsMovies() bool    { return false }
func (p *Provider) SupportsTVShows() bool   { return true }
func (p *Provider) SupportsPeople() bool    { return false }
func (p *Provider) ClearCache()             { p.client.clearCache() }

// SupportsLanguage returns true for Japanese, English, and romaji content.
func (p *Provider) SupportsLanguage(lang string) bool {
	switch lang {
	case "ja", "en", "x-jat": // x-jat = romaji
		return true
	default:
		return false
	}
}

// --- TVShowProvider ---

func (p *Provider) SearchTVShow(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	limit := 20
	matches, err := p.client.SearchAnime(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, metadata.ErrNotFound
	}

	results := make([]metadata.TVShowSearchResult, 0, len(matches))
	for _, m := range matches {
		results = append(results, mapTitleToTVShowSearchResult(m))
	}
	return results, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}

	m := mapAnimeToTVShowMetadata(anime)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}

func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}

	credits := mapCredits(anime)
	if credits == nil || (len(credits.Cast) == 0 && len(credits.Crew) == 0) {
		return nil, metadata.ErrNotFound
	}
	return credits, nil
}

func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}

	images := mapImages(anime)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}

	ids := findExternalIDs(anime)
	if ids == nil {
		return nil, metadata.ErrNotFound
	}
	return ids, nil
}

func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, _ string) (*metadata.SeasonMetadata, error) {
	aid, err := strconv.Atoi(showID)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", showID, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}
	if anime == nil {
		return nil, metadata.ErrNotFound
	}

	episodes := mapEpisodes(anime, seasonNum)
	if len(episodes) == 0 {
		return nil, metadata.ErrNotFound
	}

	sm := &metadata.SeasonMetadata{
		ProviderID:   showID,
		Provider:     metadata.ProviderAniDB,
		ShowID:       showID,
		SeasonNumber: seasonNum,
		Name:         fmt.Sprintf("Season %d", seasonNum),
		Episodes:     episodes,
	}

	if anime.Picture != "" {
		img := ImageBaseURL + anime.Picture
		sm.PosterPath = &img
	}

	return sm, nil
}

func (p *Provider) GetEpisode(ctx context.Context, showID string, _, episodeNum int, _ string) (*metadata.EpisodeMetadata, error) {
	aid, err := strconv.Atoi(showID)
	if err != nil {
		return nil, fmt.Errorf("anidb: invalid anime ID %q: %w", showID, err)
	}

	anime, err := p.client.GetAnime(ctx, aid)
	if err != nil {
		return nil, err
	}
	if anime == nil {
		return nil, metadata.ErrNotFound
	}

	for _, ep := range anime.Episodes.Episode {
		if ep.EpNo.Type != 1 { // Regular episodes only
			continue
		}
		num, err := strconv.Atoi(ep.EpNo.Text)
		if err != nil {
			continue
		}
		if num == episodeNum {
			return mapEpisodeToMetadata(ep, showID), nil
		}
	}

	return nil, metadata.ErrNotFound
}
