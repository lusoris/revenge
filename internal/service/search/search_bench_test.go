package search

import (
	"testing"
)

func BenchmarkDefaultSearchParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultSearchParams()
	}
}

func BenchmarkDefaultTVShowSearchParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultTVShowSearchParams()
	}
}

func BenchmarkSearchCacheKey(b *testing.B) {
	svc := &CachedMovieSearchService{}

	params := SearchParams{
		Query:    "The Matrix",
		Page:     1,
		PerPage:  20,
		SortBy:   "popularity:desc",
		FilterBy: "genres:=action && year:>=2000",
		FacetBy:  []string{"genres", "year", "status"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.searchCacheKey(params)
	}
}

func BenchmarkSearchCacheKey_Long(b *testing.B) {
	svc := &CachedMovieSearchService{}

	params := SearchParams{
		Query:    "extremely long search query that a user might type with lots of keywords and details about the movie they are looking for",
		Page:     3,
		PerPage:  50,
		SortBy:   "release_date:desc",
		FilterBy: "genres:=[action, thriller, sci-fi] && year:>=1990 && year:<=2025 && has_file:=true && resolution:=[1080p, 4k]",
		FacetBy:  []string{"genres", "year", "status", "directors", "resolution", "has_file"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.searchCacheKey(params)
	}
}

func BenchmarkSearchCacheKey_Parallel(b *testing.B) {
	svc := &CachedMovieSearchService{}

	params := SearchParams{
		Query:    "Breaking Bad",
		Page:     1,
		PerPage:  20,
		SortBy:   "popularity:desc",
		FilterBy: "genres:=drama",
		FacetBy:  []string{"genres", "year"},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = svc.searchCacheKey(params)
		}
	})
}
