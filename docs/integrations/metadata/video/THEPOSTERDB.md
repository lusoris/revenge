# ThePosterDB Integration

> Curated high-quality posters for movies and TV shows

**Status**: 🟡 PLANNED
**Priority**: 🟢 LOW (Phase 8 - Media Enhancements)
**Type**: HTTP scraper (no official API)

---

## Overview

ThePosterDB is a community-driven platform for high-quality, curated movie and TV show posters. Unlike TMDb (which hosts user-uploaded posters), ThePosterDB focuses on:
- **Professional-grade posters**: Textless, minimal designs
- **Set collections**: 4K logos, IMAX editions, streaming service themes
- **Consistent style**: Uniform aspect ratios and quality

**Integration Points**:
- **Web scraping**: No official API (use community API OR web scraping)
- **Poster downloads**: High-resolution poster images
- **Set collections**: Download entire poster sets (e.g., Marvel Cinematic Universe with matching style)

**⚠️ Important**: ThePosterDB has NO official API. Use unofficial community API OR web scraping with strict rate limiting.

---

## Developer Resources

- 🔗 **Website**: https://theposterdb.com/
- ❌ **No Official API**
- 🔗 **Community API**: https://github.com/jarulsamy/ThePosterDB-API (unofficial, community-maintained)
- 🔗 **robots.txt**: https://theposterdb.com/robots.txt (respect crawling rules)

---

## Integration Options

### Option 1: Community API (Recommended)

Use unofficial community-maintained API:

**GitHub**: https://github.com/jarulsamy/ThePosterDB-API

**Features**:
- Search posters by title/IMDb ID
- Download poster images
- Browse sets/collections
- User authentication (optional)

**Limitations**:
- Not officially supported by ThePosterDB
- May break if website structure changes
- Rate limiting required

### Option 2: Web Scraping (Fallback)

Scrape ThePosterDB website directly:

**Constraints**:
- Respect `robots.txt`
- Aggressive rate limiting (1 request per 5-10 seconds)
- User-Agent identification
- Handle CAPTCHAs gracefully (fail silently)

---

## Implementation Checklist

- [ ] **Community API Client** (`internal/service/metadata/provider_posterdb.go`)
  - [ ] Search posters by IMDb ID
  - [ ] Search posters by title
  - [ ] Download poster images
  - [ ] Browse sets/collections
  - [ ] Rate limiting (1 req/5s minimum)
  - [ ] Error handling (graceful degradation)

- [ ] **Poster Selection UI**
  - [ ] Display ThePosterDB posters as alternatives
  - [ ] Allow users to choose preferred poster
  - [ ] Show poster sets (e.g., MCU collection)
  - [ ] Filter by style (textless, minimal, etc.)

- [ ] **Poster Storage**
  - [ ] Download high-resolution posters
  - [ ] Store locally (configurable path)
  - [ ] Generate Blurhash for placeholders
  - [ ] Image optimization (WebP conversion)

- [ ] **Fallback Strategy**
  - [ ] Use TMDb posters as default
  - [ ] Offer ThePosterDB as alternative (optional upgrade)
  - [ ] Graceful failure (if ThePosterDB unavailable)

---

## Revenge Integration Pattern

```
User views movie details (The Matrix)
           ↓
Display TMDb poster (default)
           ↓
Show "Browse alternative posters" button
           ↓
User clicks button
           ↓
Query ThePosterDB (community API OR scraper)
           ↓
Display curated poster options
           ↓
User selects preferred poster
           ↓
Download poster from ThePosterDB
           ↓
Store locally + update PostgreSQL (movies.poster_path)
           ↓
Display new poster in UI
```

### Go Client Example (Community API)

```go
type PosterDBClient struct {
    baseURL string  // Community API base URL
    client  *http.Client
    limiter *rate.Limiter  // 1 req/5s
}

func (c *PosterDBClient) SearchByIMDbID(ctx context.Context, imdbID string) ([]Poster, error) {
    c.limiter.Wait(ctx)  // Rate limiting
    
    url := fmt.Sprintf("%s/posters?imdb=%s", c.baseURL, imdbID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("User-Agent", "Revenge Media Server/1.0 (admin@example.com)")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }
    defer resp.Body.Close()
    
    var posters []Poster
    json.NewDecoder(resp.Body).Decode(&posters)
    return posters, nil
}

func (c *PosterDBClient) DownloadPoster(ctx context.Context, posterURL string) ([]byte, error) {
    c.limiter.Wait(ctx)
    
    req, _ := http.NewRequestWithContext(ctx, "GET", posterURL, nil)
    req.Header.Set("User-Agent", "Revenge Media Server/1.0 (admin@example.com)")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("download failed: %w", err)
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}
```

---

## Related Documentation

- [Movie Module](../../architecture/modules/MOVIE.md)
- [TV Show Module](../../architecture/modules/TVSHOW.md)
- [TMDb Integration](TMDB.md) - Default poster source
- [Media Enhancements](../../features/MEDIA_ENHANCEMENTS.md)

---

## Poster Sets & Collections

ThePosterDB offers curated poster sets with matching styles:

| Set Type | Description |
|----------|-------------|
| **Textless** | Posters without text/logos |
| **Minimal** | Clean, minimalist designs |
| **4K** | 4K Ultra HD branding |
| **IMAX** | IMAX edition branding |
| **Streaming** | Netflix, Disney+, Prime Video themes |
| **MCU** | Marvel Cinematic Universe (matching style) |
| **DC** | DC Extended Universe (matching style) |
| **Star Wars** | Star Wars saga (matching style) |

---

## Notes

- **NO official API** - use community API OR web scraping
- **Community API is unofficial** - may break if ThePosterDB changes structure
- **Rate limiting CRITICAL** - minimum 1 request per 5 seconds (avoid IP ban)
- **Respect robots.txt** - https://theposterdb.com/robots.txt
- **User-Agent required** - identify Revenge Media Server (include contact email)
- **Graceful degradation** - use TMDb posters if ThePosterDB unavailable
- **Optional feature** - ThePosterDB posters are enhancement, not requirement
- **High resolution** - ThePosterDB posters are higher quality than TMDb user uploads
- **Curated content** - ThePosterDB moderators approve posters (consistent quality)
- **Set collections** - Download entire poster sets for consistency (e.g., MCU with matching style)
- **CAPTCHA handling** - If CAPTCHA detected, fail silently and use TMDb fallback
- **Legal considerations** - ThePosterDB allows personal use (check ToS for commercial use)
- **Storage requirements** - High-resolution posters require more disk space (consider compression)
