// Package circuitbreaker provides circuit breaker protection for external API clients.
//
// It wraps [github.com/sony/gobreaker] with pre-configured tiers for different
// classes of external services (public APIs, local services, CDNs) and integrates
// with the observability package for Prometheus metric exposition.
//
// Usage with imroc/req v3 clients:
//
//	client := req.C().SetBaseURL(baseURL)...
//	circuitbreaker.WrapReqClient(client, "tmdb", circuitbreaker.TierExternal)
package circuitbreaker

import (
	"log/slog"
	"time"

	"github.com/imroc/req/v3"
	"github.com/sony/gobreaker"

	"github.com/lusoris/revenge/internal/infra/observability"
)

// Tier defines the class of external service for tuning circuit breaker parameters.
type Tier int

const (
	// TierExternal is for public external APIs (TMDb, OMDb, Trakt, TVDb, etc.).
	// Tolerates more failures before tripping since external APIs may have
	// transient issues and requests already have retry+backoff.
	TierExternal Tier = iota

	// TierLocal is for local network services (Radarr, Sonarr).
	// Less tolerant — if a local service is down, it is likely fully down.
	TierLocal

	// TierCDN is for CDN/image download endpoints (image.tmdb.org, fanart.tv CDN).
	// Most tolerant — CDNs are very rarely down, failures are usually transient.
	TierCDN
)

// String returns the tier name for logging.
func (t Tier) String() string {
	switch t {
	case TierExternal:
		return "external"
	case TierLocal:
		return "local"
	case TierCDN:
		return "cdn"
	default:
		return "unknown"
	}
}

// settings returns [gobreaker.Settings] tuned for the given service tier.
func settings(name string, tier Tier) gobreaker.Settings {
	s := gobreaker.Settings{
		Name: name,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Warn("circuit breaker state change",
				slog.String("breaker", name),
				slog.String("from", from.String()),
				slog.String("to", to.String()),
			)
			observability.RecordCircuitBreakerStateChange(name, to.String())
		},
	}

	switch tier {
	case TierExternal:
		s.MaxRequests = 3
		s.Interval = 60 * time.Second
		s.Timeout = 30 * time.Second
		s.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		}
	case TierLocal:
		s.MaxRequests = 2
		s.Interval = 30 * time.Second
		s.Timeout = 15 * time.Second
		s.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		}
	case TierCDN:
		s.MaxRequests = 5
		s.Interval = 120 * time.Second
		s.Timeout = 60 * time.Second
		s.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 10
		}
	}

	return s
}

// Breaker wraps a [gobreaker.TwoStepCircuitBreaker] for use as request middleware.
type Breaker struct {
	cb *gobreaker.TwoStepCircuitBreaker
}

// New creates a circuit breaker for the given service name and tier.
func New(name string, tier Tier) *Breaker {
	return &Breaker{
		cb: gobreaker.NewTwoStepCircuitBreaker(settings(name, tier)),
	}
}

// State returns the current circuit breaker state (Closed, HalfOpen, or Open).
func (b *Breaker) State() gobreaker.State {
	return b.cb.State()
}

// Name returns the circuit breaker name.
func (b *Breaker) Name() string {
	return b.cb.Name()
}

// WrapReqClient wraps a [req.Client] with circuit breaker protection using
// req WrapRoundTripFunc middleware. This must be called after all other
// client configuration (SetProxyURL, etc.) is complete.
//
// The circuit breaker counts transport errors and HTTP 5xx responses as
// failures. HTTP 4xx and lower status codes count as successes (client
// errors are not indicative of upstream service failure).
func WrapReqClient(client *req.Client, name string, tier Tier) *Breaker {
	b := New(name, tier)

	client.WrapRoundTripFunc(func(rt req.RoundTripper) req.RoundTripFunc {
		return func(r *req.Request) (*req.Response, error) {
			done, err := b.cb.Allow()
			if err != nil {
				return nil, err
			}

			resp, rtErr := rt.RoundTrip(r)
			if rtErr != nil {
				done(false)
				return nil, rtErr
			}

			// 5xx = service failure, everything else = success
			if resp.Response != nil {
				done(resp.StatusCode < 500)
			} else {
				done(false)
			}

			return resp, nil
		}
	})

	return b
}
