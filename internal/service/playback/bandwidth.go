// Package playback provides bandwidth monitoring for external clients.
package playback

import (
	"sync"
	"time"
)

const (
	// BandwidthSampleWindow is how many samples to keep for averaging.
	BandwidthSampleWindow = 10
	// BandwidthSampleInterval is the minimum time between samples.
	BandwidthSampleInterval = 2 * time.Second
	// SafetyMargin is applied to bandwidth estimates (80% of measured).
	SafetyMargin = 0.8
)

// BandwidthSample represents a single bandwidth measurement.
type BandwidthSample struct {
	Timestamp time.Time
	Kbps      int
	Latency   time.Duration
	BytesSent int64
	Duration  time.Duration
}

// BandwidthMonitor tracks bandwidth for a client session.
type BandwidthMonitor struct {
	mu           sync.RWMutex
	samples      []BandwidthSample
	lastSampleAt time.Time
	isExternal   bool
}

// NewBandwidthMonitor creates a new bandwidth monitor.
func NewBandwidthMonitor(isExternal bool) *BandwidthMonitor {
	return &BandwidthMonitor{
		samples:    make([]BandwidthSample, 0, BandwidthSampleWindow),
		isExternal: isExternal,
	}
}

// AddSample records a bandwidth measurement.
func (m *BandwidthMonitor) AddSample(bytesSent int64, duration time.Duration, latency time.Duration) {
	if !m.isExternal {
		return // Don't track local clients
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// Rate limit samples
	if now.Sub(m.lastSampleAt) < BandwidthSampleInterval {
		return
	}
	m.lastSampleAt = now

	// Calculate kbps
	seconds := duration.Seconds()
	if seconds <= 0 {
		return
	}
	kbps := int(float64(bytesSent*8) / seconds / 1000)

	sample := BandwidthSample{
		Timestamp: now,
		Kbps:      kbps,
		Latency:   latency,
		BytesSent: bytesSent,
		Duration:  duration,
	}

	// Sliding window
	if len(m.samples) >= BandwidthSampleWindow {
		m.samples = m.samples[1:]
	}
	m.samples = append(m.samples, sample)
}

// GetEstimate returns the current bandwidth estimate with jitter.
func (m *BandwidthMonitor) GetEstimate() BandwidthEstimate {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.samples) == 0 {
		return BandwidthEstimate{
			IsReliable: false,
		}
	}

	// Calculate average and jitter
	var sum int
	minKbps := m.samples[0].Kbps
	maxKbps := m.samples[0].Kbps

	for _, s := range m.samples {
		sum += s.Kbps
		if s.Kbps < minKbps {
			minKbps = s.Kbps
		}
		if s.Kbps > maxKbps {
			maxKbps = s.Kbps
		}
	}

	avg := sum / len(m.samples)
	jitter := maxKbps - minKbps

	// Calculate variance for stability assessment
	var varianceSum int
	for _, s := range m.samples {
		diff := s.Kbps - avg
		varianceSum += diff * diff
	}
	variance := varianceSum / len(m.samples)

	// Average latency
	var latencySum time.Duration
	for _, s := range m.samples {
		latencySum += s.Latency
	}
	avgLatency := latencySum / time.Duration(len(m.samples))

	return BandwidthEstimate{
		AverageKbps:     avg,
		MinKbps:         minKbps,
		MaxKbps:         maxKbps,
		JitterKbps:      jitter,
		Variance:        variance,
		AverageLatency:  avgLatency,
		SampleCount:     len(m.samples),
		IsReliable:      len(m.samples) >= 3,
		RecommendedKbps: m.calculateRecommended(avg, jitter),
	}
}

// calculateRecommended returns a safe bitrate recommendation.
func (m *BandwidthMonitor) calculateRecommended(avgKbps, jitterKbps int) int {
	// Apply safety margin and subtract jitter
	safe := int(float64(avgKbps)*SafetyMargin) - jitterKbps
	if safe < 500 {
		safe = 500 // Minimum 500 kbps
	}
	return safe
}

// BandwidthEstimate represents the current bandwidth assessment.
type BandwidthEstimate struct {
	AverageKbps     int           `json:"average_kbps"`
	MinKbps         int           `json:"min_kbps"`
	MaxKbps         int           `json:"max_kbps"`
	JitterKbps      int           `json:"jitter_kbps"` // max - min
	Variance        int           `json:"variance"`    // Statistical variance
	AverageLatency  time.Duration `json:"average_latency"`
	SampleCount     int           `json:"sample_count"`
	IsReliable      bool          `json:"is_reliable"`      // Have enough samples
	RecommendedKbps int           `json:"recommended_kbps"` // Safe bitrate to use
}

// QualityLevel returns a suggested quality level based on bandwidth.
func (e *BandwidthEstimate) QualityLevel() string {
	switch {
	case e.RecommendedKbps >= 25000:
		return "4k"
	case e.RecommendedKbps >= 10000:
		return "1080p_high"
	case e.RecommendedKbps >= 5000:
		return "1080p"
	case e.RecommendedKbps >= 3000:
		return "720p"
	case e.RecommendedKbps >= 1500:
		return "480p"
	default:
		return "360p"
	}
}

// IsStable returns true if the connection is stable.
func (e *BandwidthEstimate) IsStable() bool {
	if !e.IsReliable {
		return false
	}
	// Jitter less than 20% of average is considered stable
	return float64(e.JitterKbps)/float64(e.AverageKbps) < 0.2
}
