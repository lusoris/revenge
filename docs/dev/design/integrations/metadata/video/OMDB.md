# OMDb (Open Movie Database) Integration

> Fallback metadata provider + IMDb ratings

**Status**: 🟡 PLANNED
**Priority**: 🟡 HIGH (Phase 2 - Movie Module)
**Type**: HTTP API client

---

## Overview

OMDb is a lightweight metadata provider for movies and TV shows, primarily used for:
- **IMDb ratings**: Authoritative ratings from IMDb
- **Fallback metadata**: When TMDb/TheTVDB data is incomplete
- **Basic movie info**: Title, year, director, plot
- **Poster fallback**: Alternative poster source

**Integration Points**:
- **API client**: Query by IMDb ID or title
- **IMDb ratings**: Primary source for IMDb ratings (not available via TMDb)
- **Fallback metadata**: Use when TMDb API fails or returns incomplete data

---

## Developer Resources

- 📚 **API Docs**: https://www.omdbapi.com/
- 🔗 **API**: http://www.omdbapi.com/

---

## API Details

**Base URL**: `http://www.omdbapi.com/`
**Authentication**: API Key query parameter `?apikey={key}`
**Rate Limits**: 1,000 requests/day (free tier)
**Free Tier**: Available (requires API key registration)

### Query Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `i` | IMDb ID | `i=tt0133093` |
| `t` | Movie/series title | `t=The Matrix` |
| `y` | Year of release | `y=1999` |
| `type` | Type (movie, series, episode) | `type=movie` |
| `plot` | Plot length (short, full) | `plot=full` |
| `apikey` | API key | `apikey=your-key` |

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/?i={imdb_id}` | Get movie/series by IMDb ID |
| `/?t={title}&y={year}` | Get movie/series by title & year |
| `/?s={search}` | Search movies/series |

---

## Response Example

```json
{
  "Title": "The Matrix",
  "Year": "1999",
  "Rated": "R",
  "Released": "31 Mar 1999",
  "Runtime": "136 min",
  "Genre": "Action, Sci-Fi",
  "Director": "Lana Wachowski, Lilly Wachowski",
  "Writer": "Lilly Wachowski, Lana Wachowski",
  "Actors": "Keanu Reeves, Laurence Fishburne, Carrie-Anne Moss",
  "Plot": "When a beautiful stranger leads computer hacker Neo to a forbidding underworld...",
  "Language": "English",
  "Country": "United States, Australia",
  "Awards": "Won 4 Oscars. 42 wins & 51 nominations total",
  "Poster": "https://m.media-amazon.com/images/M/MV5BNzQzOTk3OTAtNDQ0Zi00ZTVkLWI0MTEtMDllZjNkYzNjNTc4L2ltYWdlXkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg",
  "Ratings": [
    {
      "Source": "Internet Movie Database",
      "Value": "8.7/10"
    },
    {
      "Source": "Rotten Tomatoes",
      "Value": "83%"
    },
    {
      "Source": "Metacritic",
      "Value": "73/100"
    }
  ],
  "Metascore": "73",
  "imdbRating": "8.7",
  "imdbVotes": "1,958,348",
  "imdbID": "tt0133093",
  "Type": "movie",
  "DVD": "21 Sep 1999",
  "BoxOffice": "$172,076,928",
  "Production": "N/A",
  "Website": "N/A",
  "Response": "True"
}
```

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/metadata/provider_omdb.go`)
  - [ ] Query by IMDb ID
  - [ ] Query by title + year
  - [ ] Search functionality
  - [ ] Rate limiting (1,000 req/day)
  - [ ] Error handling & retries

- [ ] **IMDb Ratings**
  - [ ] Extract `imdbRating` field
  - [ ] Extract `imdbVotes` field
  - [ ] Store in Revenge `movies` table (`imdb_rating`, `imdb_votes`)
  - [ ] Display IMDb rating badge in UI

- [ ] **Fallback Metadata**
  - [ ] Use when TMDb API fails
  - [ ] Use for missing fields (director, plot, etc.)
  - [ ] Poster fallback (if TMDb poster unavailable)

- [ ] **Ratings Aggregation**
  - [ ] IMDb rating (primary)
  - [ ] Rotten Tomatoes rating (secondary)
  - [ ] Metacritic score (secondary)

---

## Revenge Integration Pattern

```
Fetch movie metadata from TMDb
           ↓
TMDb returns metadata (without IMDb rating)
           ↓
Query OMDb by IMDb ID for rating
           ↓
OMDb returns IMDb rating (8.7/10)
           ↓
Store IMDb rating in PostgreSQL (movies.imdb_rating)
           ↓
Display IMDb rating in UI
```

### Go Client Example

```go
type OMDbClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
    limiter *rate.Limiter  // 1,000 req/day
}

func (c *OMDbClient) GetByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
    c.limiter.Wait(ctx)  // Rate limiting

    url := fmt.Sprintf("%s?i=%s&apikey=%s&plot=full", c.baseURL, imdbID, c.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get movie: %w", err)
    }
    defer resp.Body.Close()

    var movie Movie
    json.NewDecoder(resp.Body).Decode(&movie)

    if movie.Response == "False" {
        return nil, fmt.Errorf("movie not found: %s", movie.Error)
    }

    return &movie, nil
}

func (c *OMDbClient) GetByTitle(ctx context.Context, title string, year int) (*Movie, error) {
    c.limiter.Wait(ctx)

    url := fmt.Sprintf("%s?t=%s&y=%d&apikey=%s&plot=full",
        c.baseURL, url.QueryEscape(title), year, c.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get movie: %w", err)
    }
    defer resp.Body.Close()

    var movie Movie
    json.NewDecoder(resp.Body).Decode(&movie)
    return &movie, nil
}
```

---

## Related Documentation

- [Movie Module](../../architecture/modules/MOVIE.md)
- [TMDb Integration](TMDB.md) - Primary movie metadata source
- [Metadata Enrichment Pattern](../../patterns/metadata_enrichment.md)

---

## Rating Sources

OMDb aggregates ratings from multiple sources:

| Source | Field | Format |
|--------|-------|--------|
| IMDb | `imdbRating` | `8.7/10` |
| Rotten Tomatoes | `Ratings[1].Value` | `83%` |
| Metacritic | `Metascore` | `73/100` |

**Recommended**: Display IMDb rating prominently (most widely recognized).

---

## Notes

- **OMDb is lightweight** (simpler API than TMDb/TheTVDB)
- **IMDb ratings are authoritative** (OMDb fetches from IMDb)
- **Rate limit**: 1,000 requests/day (free tier) - use strategically for IMDb ratings only
- **Paid tiers available**: Higher rate limits (http://www.omdbapi.com/apikey.aspx)
- **Poster quality**: Lower resolution than TMDb (use TMDb posters when available)
- **Plot field**: Use `plot=full` for complete plot summary
- **Response validation**: Check `Response` field (`"True"` = success, `"False"` = error)
- **Error handling**: OMDb returns `{"Response":"False","Error":"Movie not found!"}` on error
- **No authentication complexity**: Simple API key (no JWT, no OAuth)
- **Use case**: Primary use = IMDb ratings, fallback metadata when TMDb incomplete
- **Rotten Tomatoes**: OMDb provides Rotten Tomatoes scores (not available via TMDb)
- **Box office**: OMDb provides box office revenue (useful metric)
- **Awards**: OMDb provides Oscar/award info (not available via TMDb)
