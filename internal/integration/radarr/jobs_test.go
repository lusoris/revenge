package radarr

import (
	"context"
	"testing"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadarrSyncJobArgs_Kind(t *testing.T) {
	t.Parallel()

	args := RadarrSyncJobArgs{}
	assert.Equal(t, "radarr_sync", args.Kind())
}

func TestRadarrSyncJobArgs_Operations(t *testing.T) {
	t.Parallel()

	assert.Equal(t, RadarrSyncOperation("full"), RadarrSyncOperationFull)
	assert.Equal(t, RadarrSyncOperation("single"), RadarrSyncOperationSingle)
}

func TestNewRadarrSyncWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewRadarrSyncWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.syncService)
	assert.NotNil(t, worker.logger)
}

func TestRadarrSyncWorker_Work_NilService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewRadarrSyncWorker(nil, logger)

	job := &river.Job[RadarrSyncJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "radarr_sync"},
		Args: RadarrSyncJobArgs{
			Operation: RadarrSyncOperationFull,
		},
	}

	// Should return nil when syncService is nil
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestRadarrSyncWorker_Work_UnknownOperation(t *testing.T) {
	t.Parallel()

	// Need a non-nil syncService to test the operation switch
	// For this test, we'll skip since it requires a full syncService
	t.Skip("Requires mock syncService")
}

func TestRadarrWebhookJobArgs_Kind(t *testing.T) {
	t.Parallel()

	args := RadarrWebhookJobArgs{}
	assert.Equal(t, "radarr_webhook", args.Kind())
}

func TestNewRadarrWebhookWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewRadarrWebhookWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.webhookHandler)
	assert.NotNil(t, worker.logger)
}

func TestRadarrWebhookWorker_Work_NilHandler(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewRadarrWebhookWorker(nil, logger)

	job := &river.Job[RadarrWebhookJobArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "radarr_webhook"},
		Args: RadarrWebhookJobArgs{
			Payload: WebhookPayload{
				EventType: "Test",
				Movie:     &WebhookMovie{ID: 123},
			},
		},
	}

	// Should return nil when webhookHandler is nil
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}
