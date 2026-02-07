package sonarr

import (
	"context"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

func TestSonarrSyncJobArgs_Kind(t *testing.T) {
	t.Parallel()

	args := SonarrSyncJobArgs{}
	assert.Equal(t, "sonarr_sync", args.Kind())
}

func TestSonarrSyncJobArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := SonarrSyncJobArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueHigh, opts.Queue)
}

func TestSonarrSyncJobArgs_Operations(t *testing.T) {
	t.Parallel()

	assert.Equal(t, SonarrSyncOperation("full"), SonarrSyncOperationFull)
	assert.Equal(t, SonarrSyncOperation("single"), SonarrSyncOperationSingle)
}

func TestNewSonarrSyncWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrSyncWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.syncService)
	assert.NotNil(t, worker.logger)
}

func TestSonarrSyncWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrSyncWorker(nil, logger)

	job := &river.Job[SonarrSyncJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "sonarr_sync"},
		Args:   SonarrSyncJobArgs{Operation: SonarrSyncOperationFull},
	}

	assert.Equal(t, 10*time.Minute, worker.Timeout(job))
}

func TestSonarrSyncWorker_Work_NilService(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrSyncWorker(nil, logger)

	job := &river.Job[SonarrSyncJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "sonarr_sync"},
		Args: SonarrSyncJobArgs{
			Operation: SonarrSyncOperationFull,
		},
	}

	// Should return nil when syncService is nil
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestSonarrWebhookJobArgs_Kind(t *testing.T) {
	t.Parallel()

	args := SonarrWebhookJobArgs{}
	assert.Equal(t, "sonarr_webhook", args.Kind())
}

func TestSonarrWebhookJobArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := SonarrWebhookJobArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueHigh, opts.Queue)
}

func TestNewSonarrWebhookWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrWebhookWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.webhookHandler)
	assert.NotNil(t, worker.logger)
}

func TestSonarrWebhookWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrWebhookWorker(nil, logger)

	job := &river.Job[SonarrWebhookJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "sonarr_webhook"},
		Args: SonarrWebhookJobArgs{
			Payload: WebhookPayload{EventType: "Test"},
		},
	}

	assert.Equal(t, 1*time.Minute, worker.Timeout(job))
}

func TestSonarrWebhookWorker_Work_NilHandler(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewSonarrWebhookWorker(nil, logger)

	job := &river.Job[SonarrWebhookJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "sonarr_webhook"},
		Args: SonarrWebhookJobArgs{
			Payload: WebhookPayload{
				EventType: "Test",
				Series:    &WebhookSeries{ID: 123},
			},
		},
	}

	// Should return nil when webhookHandler is nil
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}
