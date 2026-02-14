package circuitbreaker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/imroc/req/v3"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tier Tier
	}{
		{"external-api", TierExternal},
		{"local-service", TierLocal},
		{"cdn-service", TierCDN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(tt.name, tt.tier)
			require.NotNil(t, b)
			assert.Equal(t, tt.name, b.Name())
			assert.Equal(t, gobreaker.StateClosed, b.State())
		})
	}
}

func TestTierString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "external", TierExternal.String())
	assert.Equal(t, "local", TierLocal.String())
	assert.Equal(t, "cdn", TierCDN.String())
	assert.Equal(t, "unknown", Tier(99).String())
}

func TestWrapReqClient_Success(t *testing.T) {
	t.Parallel()

	var callCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := req.C().SetBaseURL(srv.URL).SetCommonRetryCount(0)
	b := WrapReqClient(client, "test-success", TierExternal)

	resp, err := client.R().Get("/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, gobreaker.StateClosed, b.State())
	assert.Equal(t, int32(1), callCount.Load())
}

func TestWrapReqClient_ServerErrors_TripBreaker(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := req.C().
		SetBaseURL(srv.URL).
		DisableAutoReadResponse().
		SetCommonRetryCount(0)

	// TierLocal trips after 3 consecutive failures
	b := WrapReqClient(client, "test-trip", TierLocal)

	for i := 0; i < 3; i++ {
		resp, err := client.R().Get("/fail")
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	}

	assert.Equal(t, gobreaker.StateOpen, b.State())

	// Next request should fail without hitting the server
	_, err := client.R().Get("/fail")
	require.Error(t, err)
	assert.True(t, errors.Is(err, gobreaker.ErrOpenState))
}

func TestWrapReqClient_ClientErrors_DontTrip(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	client := req.C().
		SetBaseURL(srv.URL).
		DisableAutoReadResponse().
		SetCommonRetryCount(0)

	b := WrapReqClient(client, "test-4xx", TierLocal)

	for i := 0; i < 10; i++ {
		resp, err := client.R().Get("/bad")
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	assert.Equal(t, gobreaker.StateClosed, b.State())
}

func TestWrapReqClient_Recovery(t *testing.T) {
	t.Parallel()

	var shouldFail atomic.Bool
	shouldFail.Store(true)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if shouldFail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := req.C().
		SetBaseURL(srv.URL).
		DisableAutoReadResponse().
		SetCommonRetryCount(0)

	// Use custom settings with very short timeout for test
	b := &Breaker{
		cb: gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
			Name:        "test-recovery",
			MaxRequests: 1,
			Timeout:     100 * time.Millisecond,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= 2
			},
		}),
	}

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
			if resp.Response != nil {
				done(resp.StatusCode < 500)
			} else {
				done(false)
			}
			return resp, nil
		}
	})

	// Trip the breaker
	for i := 0; i < 2; i++ {
		client.R().Get("/fail")
	}
	assert.Equal(t, gobreaker.StateOpen, b.State())

	// Wait for timeout -> half-open
	time.Sleep(150 * time.Millisecond)

	// Fix the server
	shouldFail.Store(false)

	// Should succeed and close the breaker
	resp, err := client.R().Get("/ok")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, gobreaker.StateClosed, b.State())
}
