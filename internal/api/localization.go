package api

import (
	"context"
	"strings"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
)

// GetMetadataLanguage returns the user's preferred metadata language.
// Priority: User preference (metadata_language) > Accept-Language header > Default (en)
func (h *Handler) GetMetadataLanguage(ctx context.Context) string {
	// 1. Try user's metadata_language preference
	userID, err := GetUserID(ctx)
	if err == nil {
		prefs, err := h.userService.GetUserPreferences(ctx, userID)
		if err == nil && prefs.MetadataLanguage != nil && *prefs.MetadataLanguage != "" {
			return *prefs.MetadataLanguage
		}
	}

	// 2. Try Accept-Language header from request metadata
	meta := middleware.GetRequestMetadata(ctx)
	if lang := parseAcceptLanguage(meta.AcceptLanguage); lang != "" {
		return lang
	}

	// 3. Default to English
	return "en"
}

// parseAcceptLanguage parses Accept-Language header and returns best match.
// Example: "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7" → "de"
func parseAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}

	// Split by comma to get language preferences
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}

	// Take first preference (highest priority)
	firstLang := strings.TrimSpace(parts[0])

	// Remove quality value if present (e.g., "de-DE;q=0.9" → "de-DE")
	if idx := strings.Index(firstLang, ";"); idx != -1 {
		firstLang = firstLang[:idx]
	}

	// Extract language code (e.g., "de-DE" → "de")
	if idx := strings.Index(firstLang, "-"); idx != -1 {
		return firstLang[:idx]
	}

	return firstLang
}

// LocalizeMovie returns a localized copy of the movie with fields in the preferred language.
func LocalizeMovie(m *movie.Movie, lang string) *movie.Movie {
	if m == nil {
		return nil
	}

	localized := *m
	localized.Title = m.GetTitle(lang)

	if overview := m.GetOverview(lang); overview != "" {
		localized.Overview = &overview
	}

	if tagline := m.GetTagline(lang); tagline != "" {
		localized.Tagline = &tagline
	}

	return &localized
}

// LocalizeMovies localizes a slice of movies.
func LocalizeMovies(movies []movie.Movie, lang string) []movie.Movie {
	localized := make([]movie.Movie, len(movies))
	for i := range movies {
		localized[i] = *LocalizeMovie(&movies[i], lang)
	}
	return localized
}

// LocalizeSeries returns a localized copy of the series with fields in the preferred language.
func LocalizeSeries(s *tvshow.Series, lang string) *tvshow.Series {
	if s == nil {
		return nil
	}

	localized := *s
	localized.Title = s.GetTitle(lang)

	if overview := s.GetOverview(lang); overview != "" {
		localized.Overview = &overview
	}

	if tagline := s.GetTagline(lang); tagline != "" {
		localized.Tagline = &tagline
	}

	return &localized
}

// LocalizeSeriesList localizes a slice of series.
func LocalizeSeriesList(series []tvshow.Series, lang string) []tvshow.Series {
	localized := make([]tvshow.Series, len(series))
	for i := range series {
		localized[i] = *LocalizeSeries(&series[i], lang)
	}
	return localized
}

// LocalizeContinueWatchingItem localizes a continue watching item.
func LocalizeContinueWatchingItem(item *movie.ContinueWatchingItem, lang string) *movie.ContinueWatchingItem {
	if item == nil {
		return nil
	}

	localized := *item

	tempMovie := movie.Movie{
		Title:         item.Title,
		Overview:      item.Overview,
		Tagline:       item.Tagline,
		OriginalTitle: item.OriginalTitle,
		TitlesI18n:    item.TitlesI18n,
		TaglinesI18n:  item.TaglinesI18n,
		OverviewsI18n: item.OverviewsI18n,
	}

	localized.Title = tempMovie.GetTitle(lang)

	if overview := tempMovie.GetOverview(lang); overview != "" {
		localized.Overview = &overview
	}

	if tagline := tempMovie.GetTagline(lang); tagline != "" {
		localized.Tagline = &tagline
	}

	return &localized
}

// LocalizeContinueWatchingItems localizes a slice of continue watching items.
func LocalizeContinueWatchingItems(items []movie.ContinueWatchingItem, lang string) []movie.ContinueWatchingItem {
	localized := make([]movie.ContinueWatchingItem, len(items))
	for i := range items {
		localized[i] = *LocalizeContinueWatchingItem(&items[i], lang)
	}
	return localized
}

// LocalizeWatchedMovieItem localizes a watched movie item.
func LocalizeWatchedMovieItem(item *movie.WatchedMovieItem, lang string) *movie.WatchedMovieItem {
	if item == nil {
		return nil
	}

	localized := *item

	tempMovie := movie.Movie{
		Title:         item.Title,
		Overview:      item.Overview,
		Tagline:       item.Tagline,
		OriginalTitle: item.OriginalTitle,
		TitlesI18n:    item.TitlesI18n,
		TaglinesI18n:  item.TaglinesI18n,
		OverviewsI18n: item.OverviewsI18n,
	}

	localized.Title = tempMovie.GetTitle(lang)

	if overview := tempMovie.GetOverview(lang); overview != "" {
		localized.Overview = &overview
	}

	if tagline := tempMovie.GetTagline(lang); tagline != "" {
		localized.Tagline = &tagline
	}

	return &localized
}

// LocalizeWatchedMovieItems localizes a slice of watched movie items.
func LocalizeWatchedMovieItems(items []movie.WatchedMovieItem, lang string) []movie.WatchedMovieItem {
	localized := make([]movie.WatchedMovieItem, len(items))
	for i := range items {
		localized[i] = *LocalizeWatchedMovieItem(&items[i], lang)
	}
	return localized
}
