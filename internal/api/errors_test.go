package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/errors"
)

func TestAPIError_Error(t *testing.T) {
	t.Run("with underlying error", func(t *testing.T) {
		baseErr := errors.New("database connection failed")
		apiErr := api.NewAPIError(http.StatusInternalServerError, "Internal error", baseErr)

		assert.Equal(t, "database connection failed", apiErr.Error())
	})

	t.Run("without underlying error", func(t *testing.T) {
		apiErr := api.NewAPIError(http.StatusNotFound, "Resource not found", nil)

		assert.Equal(t, "Resource not found", apiErr.Error())
	})
}

func TestAPIError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	apiErr := api.NewAPIError(http.StatusInternalServerError, "Internal error", baseErr)

	unwrapped := apiErr.Unwrap()
	assert.Equal(t, baseErr, unwrapped)
}

func TestNewAPIError(t *testing.T) {
	baseErr := errors.New("test error")
	apiErr := api.NewAPIError(http.StatusBadRequest, "Bad request", baseErr)

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "Bad request", apiErr.Message)
	assert.Equal(t, baseErr, apiErr.Err)
	assert.Nil(t, apiErr.Details)
}

func TestAPIError_WithDetails(t *testing.T) {
	apiErr := api.NewAPIError(http.StatusBadRequest, "Validation failed", nil)
	details := map[string]interface{}{
		"field": "email",
		"issue": "invalid format",
	}

	result := apiErr.WithDetails(details)

	require.NotNil(t, result)
	assert.Equal(t, details, result.Details)
	assert.Equal(t, apiErr, result) // Should return same instance
}

func TestToAPIError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "nil error",
			err:          nil,
			expectedCode: 0,
			expectedMsg:  "",
		},
		{
			name:         "ErrNotFound",
			err:          errors.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "Resource not found",
		},
		{
			name:         "ErrUnauthorized",
			err:          errors.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Authentication required",
		},
		{
			name:         "ErrForbidden",
			err:          errors.ErrForbidden,
			expectedCode: http.StatusForbidden,
			expectedMsg:  "Access forbidden",
		},
		{
			name:         "ErrConflict",
			err:          errors.ErrConflict,
			expectedCode: http.StatusConflict,
			expectedMsg:  "Resource conflict",
		},
		{
			name:         "ErrValidation",
			err:          errors.ErrValidation,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Validation failed",
		},
		{
			name:         "ErrBadRequest",
			err:          errors.ErrBadRequest,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Bad request",
		},
		{
			name:         "ErrUnavailable",
			err:          errors.ErrUnavailable,
			expectedCode: http.StatusServiceUnavailable,
			expectedMsg:  "Service unavailable",
		},
		{
			name:         "ErrTimeout",
			err:          errors.ErrTimeout,
			expectedCode: http.StatusGatewayTimeout,
			expectedMsg:  "Request timeout",
		},
		{
			name:         "unknown error",
			err:          errors.New("unexpected error"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := api.ToAPIError(tt.err)

			if tt.err == nil {
				assert.Nil(t, apiErr)
				return
			}

			require.NotNil(t, apiErr)
			assert.Equal(t, tt.expectedCode, apiErr.Code)
			assert.Equal(t, tt.expectedMsg, apiErr.Message)
		})
	}
}

func TestToAPIError_WrappedErrors(t *testing.T) {
	t.Run("wrapped sentinel error", func(t *testing.T) {
		baseErr := errors.ErrNotFound
		wrappedErr := errors.Wrap(baseErr, "user not found in database")

		apiErr := api.ToAPIError(wrappedErr)

		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.Code)
		assert.Equal(t, "Resource not found", apiErr.Message)
		assert.True(t, errors.Is(apiErr.Err, errors.ErrNotFound))
	})

	t.Run("multiple levels of wrapping", func(t *testing.T) {
		baseErr := errors.ErrUnauthorized
		err1 := errors.Wrap(baseErr, "token expired")
		err2 := errors.Wrap(err1, "authentication failed")

		apiErr := api.ToAPIError(err2)

		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
		assert.True(t, errors.Is(apiErr.Err, errors.ErrUnauthorized))
	})
}

func TestNotFoundError(t *testing.T) {
	apiErr := api.NotFoundError("User not found")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, "User not found", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrNotFound))
}

func TestUnauthorizedError(t *testing.T) {
	apiErr := api.UnauthorizedError("Invalid credentials")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
	assert.Equal(t, "Invalid credentials", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrUnauthorized))
}

func TestForbiddenError(t *testing.T) {
	apiErr := api.ForbiddenError("Insufficient permissions")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Code)
	assert.Equal(t, "Insufficient permissions", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrForbidden))
}

func TestConflictError(t *testing.T) {
	apiErr := api.ConflictError("Email already exists")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusConflict, apiErr.Code)
	assert.Equal(t, "Email already exists", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrConflict))
}

func TestValidationError(t *testing.T) {
	apiErr := api.ValidationError("Invalid email format")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "Invalid email format", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrValidation))
}

func TestBadRequestError(t *testing.T) {
	apiErr := api.BadRequestError("Malformed JSON")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "Malformed JSON", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrBadRequest))
}

func TestInternalError(t *testing.T) {
	baseErr := errors.New("database failure")
	apiErr := api.InternalError("Service error", baseErr)

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
	assert.Equal(t, "Service error", apiErr.Message)
	assert.Equal(t, baseErr, apiErr.Err)
}

func TestUnavailableError(t *testing.T) {
	apiErr := api.UnavailableError("Database is down")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusServiceUnavailable, apiErr.Code)
	assert.Equal(t, "Database is down", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrUnavailable))
}

func TestTimeoutError(t *testing.T) {
	apiErr := api.TimeoutError("Request took too long")

	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusGatewayTimeout, apiErr.Code)
	assert.Equal(t, "Request took too long", apiErr.Message)
	assert.True(t, errors.Is(apiErr.Err, errors.ErrTimeout))
}

func TestAPIError_ErrorConstructorsWithDetails(t *testing.T) {
	t.Run("validation error with field details", func(t *testing.T) {
		apiErr := api.ValidationError("Validation failed").WithDetails(map[string]interface{}{
			"fields": []string{"email", "password"},
		})

		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.Code)
		assert.NotNil(t, apiErr.Details)
		assert.Contains(t, apiErr.Details, "fields")
	})

	t.Run("conflict error with duplicate details", func(t *testing.T) {
		apiErr := api.ConflictError("Duplicate entry").WithDetails(map[string]interface{}{
			"duplicate_field": "username",
			"existing_id":     12345,
		})

		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusConflict, apiErr.Code)
		assert.Equal(t, "username", apiErr.Details["duplicate_field"])
		assert.Equal(t, 12345, apiErr.Details["existing_id"])
	})
}
