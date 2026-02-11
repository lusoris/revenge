package tvmaze

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// Ensure Provider implements required interfaces.
// TVmaze is a free, no-auth TV-focused provider with good episode data,
// cast/crew info, and external ID cross-referencing. It does NOT support
// movies, people, or collections.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for TVmaze.
type Provider struct {
	metadata.TVShowProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new TVmaze provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create tvmaze provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 50, // Below TMDb (100), TVDb (80), Fanart.tv (60), above OMDb (40)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderTVmaze }
func (p *Provider) Name() string                   { return "TVmaze" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return false }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true } // English-primary
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// --- TVShowProvider ---

func (p *Provider) SearchTVShow(ctx context.Context, query string, _ metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	results, err := p.client.SearchShows(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.TVShowSearchResult, 0, len(results))
	for _, r := range results {
		searchResults = append(searchResults, mapShowToTVShowSearchResult(r.Show))
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	showID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", id, err)
	}

	show, err := p.client.GetShow(ctx, showID)
	if err != nil {
		return nil, err
	}
	m := mapShowToTVShowMetadata(show)
	if m == nil {
		return nil, metadata.ErrNotFound
	}

	// Fetch seasons to populate summary
	seasons, err := p.client.GetSeasons(ctx, showID)
	if err == nil && len(seasons) > 0 {
		m.Seasons = mapSeasons(seasons)
		m.NumberOfSeasons = len(seasons)
	}

	return m, nil
}

func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	showID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", id, err)
	}

	cast, err := p.client.GetCast(ctx, showID)
	if err != nil {
		return nil, err
	}

	crew, err := p.client.GetCrew(ctx, showID)
	if err != nil {
		return nil, err
	}

	if len(cast) == 0 && len(crew) == 0 {
		return nil, metadata.ErrNotFound
	}

	return mapCast(cast, crew), nil
}

func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	showID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", id, err)
	}

	imgs, err := p.client.GetImages(ctx, showID)
	if err != nil {
		return nil, err
	}

	images := mapImages(imgs)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}



func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	showID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", id, err)
	}

	show, err := p.client.GetShow(ctx, showID)
	if err != nil {
		return nil, err
	}

	ids := mapExternalIDs(show)
	if ids == nil {
		return nil, metadata.ErrNotFound
	}
	return ids, nil
}

func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, _ string) (*metadata.SeasonMetadata, error) {
	sid, err := strconv.Atoi(showID)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", showID, err)
	}

	// Get seasons to find season metadata
	seasons, err := p.client.GetSeasons(ctx, sid)
	if err != nil {
		return nil, err
	}

	var matchedSeason *Season
	for i := range seasons {
		if seasons[i].Number == seasonNum {
			matchedSeason = &seasons[i]
			break
		}
	}
	if matchedSeason == nil {
		return nil, metadata.ErrNotFound
	}

	sm := &metadata.SeasonMetadata{
		ProviderID:   strconv.Itoa(matchedSeason.ID),
		Provider:     metadata.ProviderTVmaze,
		ShowID:       showID,
		SeasonNumber: matchedSeason.Number,
		Name:         matchedSeason.Name,
		AirDate:      parseDate(matchedSeason.PremiereDate),
	}
	if matchedSeason.Summary != nil {
		overview := stripHTML(*matchedSeason.Summary)
		sm.Overview = &overview
	}
	if matchedSeason.Image != nil && matchedSeason.Image.Original != "" {
		sm.PosterPath = &matchedSeason.Image.Original
	}

	// Get episodes for this season
	episodes, err := p.client.GetEpisodes(ctx, sid)
	if err == nil {
		sm.Episodes = mapEpisodes(episodes, seasonNum)
	}

	return sm, nil
}



func (p *Provider) GetEpisode(ctx context.Context, showID string, seasonNum, episodeNum int, _ string) (*metadata.EpisodeMetadata, error) {
	sid, err := strconv.Atoi(showID)
	if err != nil {
		return nil, fmt.Errorf("tvmaze: invalid show ID %q: %w", showID, err)
	}

	episodes, err := p.client.GetEpisodes(ctx, sid)
	if err != nil {
		return nil, err
	}

	for _, ep := range episodes {
		if ep.Season == seasonNum && ep.Number != nil && *ep.Number == episodeNum {
			em := &metadata.EpisodeMetadata{
				ProviderID:    strconv.Itoa(ep.ID),
				Provider:      metadata.ProviderTVmaze,
				ShowID:        showID,
				SeasonNumber:  seasonNum,
				EpisodeNumber: episodeNum,
				Name:          ep.Name,
				AirDate:       parseAirdate(ep.Airdate),
			}
			if ep.Runtime != nil {
				rt := util.SafeIntToInt32(*ep.Runtime)
				em.Runtime = &rt
			}
			if ep.Rating.Average != nil {
				em.VoteAverage = *ep.Rating.Average
			}
			if ep.Summary != nil {
				overview := stripHTML(*ep.Summary)
				em.Overview = &overview
			}
			if ep.Image != nil && ep.Image.Original != "" {
				em.StillPath = &ep.Image.Original
			}
			return em, nil
		}
	}

	return nil, metadata.ErrNotFound
}


