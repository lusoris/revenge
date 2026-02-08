package observability

import (
	"testing"
)

func BenchmarkHTTPRequestsTotal_Inc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/movies", "200").Inc()
	}
}

func BenchmarkHTTPRequestDuration_Observe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HTTPRequestDuration.WithLabelValues("GET", "/api/v1/movies").Observe(0.042)
	}
}

func BenchmarkCacheHit_Record(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordCacheHit("movies", "l1")
	}
}

func BenchmarkCacheMiss_Record(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordCacheMiss("movies", "l1")
	}
}

func BenchmarkJobEnqueued_Record(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordJobEnqueued("movie_library_scan")
	}
}

func BenchmarkAuthAttempt_Record(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecordAuthAttempt("login", "success")
	}
}

func BenchmarkDBQueryDuration_Observe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DBQueryDuration.WithLabelValues("select").Observe(0.003)
	}
}

func BenchmarkSearchQueryDuration_Observe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SearchQueryDuration.WithLabelValues("search").Observe(0.015)
	}
}
