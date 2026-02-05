package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/lusoris/revenge/internal/content/movie"
)

// GetUserLanguage extracts preferred language from context.
// Priority: User settings > Accept-Language header > Default (en)
func GetUserLanguage(ctx context.Context) string {
	// 1. Try to get from user settings (from auth context)
	// TODO: Implement when user settings are available
	// if user := GetUserFromContext(ctx); user != nil && user.PreferredLanguage != "" {
	//     return user.PreferredLanguage
	// }

	// 2. Try Accept-Language header
	if req := getRequestFromContext(ctx); req != nil {
		if lang := parseAcceptLanguage(req.Header.Get("Accept-Language")); lang != "" {
			return lang
		}
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

// getRequestFromContext extracts *http.Request from context.
// This depends on the framework being used (fiber, chi, stdlib, etc.)
func getRequestFromContext(ctx context.Context) *http.Request {
	// Try to get request from context
	// Note: This is framework-specific. Adjust based on actual framework used.
	if req, ok := ctx.Value(http.Request{}).(*http.Request); ok {
		return req
	}
	return nil
}

// LocalizeMovie returns a localized copy of the movie with fields in the preferred language.
func LocalizeMovie(m *movie.Movie, lang string) *movie.Movie {
	if m == nil {
		return nil
	}

	// Create a copy to avoid modifying the original
	localized := *m

	// Override default fields with localized versions using fallback logic
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

// LocalizeContinueWatchingItem localizes a continue watching item.
func LocalizeContinueWatchingItem(item *movie.ContinueWatchingItem, lang string) *movie.ContinueWatchingItem {
	if item == nil {
		return nil
	}

	// Create a copy to avoid modifying the original
	localized := *item

	// Create a temporary Movie to use GetTitle/GetOverview/GetTagline methods
	tempMovie := movie.Movie{
		Title:         item.Title,
		Overview:      item.Overview,
		Tagline:       item.Tagline,
		OriginalTitle: item.OriginalTitle,
		TitlesI18n:    item.TitlesI18n,
		TaglinesI18n:  item.TaglinesI18n,
		OverviewsI18n: item.OverviewsI18n,
	}

	// Override default fields with localized versions
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

	// Create a copy to avoid modifying the original
	localized := *item

	// Create a temporary Movie to use GetTitle/GetOverview/GetTagline methods
	tempMovie := movie.Movie{
		Title:         item.Title,
		Overview:      item.Overview,
		Tagline:       item.Tagline,
		OriginalTitle: item.OriginalTitle,
		TitlesI18n:    item.TitlesI18n,
		TaglinesI18n:  item.TaglinesI18n,
		OverviewsI18n: item.OverviewsI18n,
	}

	// Override default fields with localized versions
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
