# Error Handling

> Error flow from database to HTTP response. Written from code as of 2026-02-06.

---

## Error Flow

```
PostgreSQL (pgx)
    |
    v
Repository Layer          pgx.ErrNoRows → domain sentinel error
    |                     other errors → fmt.Errorf("context: %w", err)
    v
Service Layer             domain logic errors → package-level sentinel vars
    |                     wraps repo errors with context
    v
Handler Layer             errors.Is() switch → ogen response type
    |                     unknown errors → 500 with generic message
    v
HTTP Response (ogen)      typed response structs (generated from OpenAPI)
```

---

## Layer 1: Repository Errors

Repositories convert database errors to domain errors. Two patterns exist:

**Pattern A: Sentinel conversion** (activity, metadata):

```go
// internal/service/activity/repository_pg.go
func (r *RepositoryPg) Get(ctx context.Context, id uuid.UUID) (*Entry, error) {
    result, err := r.queries.GetActivityLog(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound  // Package-level sentinel
        }
        return nil, err
    }
    return dbActivityToEntry(result), nil
}
```

**Pattern B: Wrapped errors** (user, auth):

```go
// internal/service/user/repository_pg.go
func (r *postgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*db.SharedUser, error) {
    user, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user by ID: %w", err)
    }
    return &user, nil
}
```

Pattern A is preferred — it lets callers use `errors.Is()` without knowing about pgx.

---

## Layer 2: Sentinel Errors

### Central sentinels (`internal/errors/errors.go`)

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrConflict     = errors.New("conflict")
    ErrValidation   = errors.New("validation failed")
    ErrBadRequest   = errors.New("bad request")
    ErrInternal     = errors.New("internal server error")
    ErrUnavailable  = errors.New("service unavailable")
    ErrTimeout      = errors.New("timeout")
)
```

Helper: `errors.Wrap(err, "context message")` wraps while preserving `errors.Is()` chain.

### Package-level sentinels

Services define their own sentinel errors for domain-specific cases:

```go
// internal/service/apikeys/service.go
var (
    ErrKeyNotFound      = errors.New("API key not found")
    ErrKeyExpired       = errors.New("API key has expired")
    ErrMaxKeysExceeded  = errors.New("maximum number of API keys exceeded")
    ErrInvalidKeyFormat = errors.New("invalid API key format")
    ErrInvalidScope     = errors.New("invalid scope")
)
```

### Custom error types

For errors that carry structured data (metadata service):

```go
// internal/service/metadata/errors.go
type ProviderError struct {
    Provider   ProviderID
    StatusCode int
    Message    string
    Err        error
}

func (e *ProviderError) Error() string {
    return fmt.Sprintf("metadata provider %s (status %d): %s", e.Provider, e.StatusCode, e.Message)
}

func (e *ProviderError) Unwrap() error { return e.Err }
```

`AggregateError` collects multiple errors (e.g., when multiple providers fail):

```go
type AggregateError struct { Errors []error }
func (e *AggregateError) Add(err error) { ... }
```

---

## Layer 3: Service Error Wrapping

Services wrap repository and internal errors with context using `fmt.Errorf`:

```go
// internal/service/auth/service.go
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
    emailToken, err := s.repo.GetEmailVerificationToken(ctx, tokenHash)
    if err != nil {
        return fmt.Errorf("invalid or expired verification token: %w", err)
    }

    if err := txQueries.MarkEmailVerificationTokenUsed(ctx, emailToken.ID); err != nil {
        return fmt.Errorf("failed to mark token as used: %w", err)
    }
    // ...
}
```

**Rules:**
- Always wrap with `%w` to preserve the error chain
- Prefix with what failed: `"failed to X: %w"`
- Never expose internal details in error messages that reach the user

---

## Layer 4: Handler Error Mapping

### APIError type (`internal/api/errors.go`)

```go
type APIError struct {
    Code    int                    `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
    Err     error                  `json:"-"`  // Never serialized
}
```

### Automatic mapping: `ToAPIError()`

```go
func ToAPIError(err error) *APIError {
    switch {
    case errors.Is(err, errors.ErrNotFound):
        return NewAPIError(http.StatusNotFound, "Resource not found", err)
    case errors.Is(err, errors.ErrUnauthorized):
        return NewAPIError(http.StatusUnauthorized, "Authentication required", err)
    case errors.Is(err, errors.ErrForbidden):
        return NewAPIError(http.StatusForbidden, "Access forbidden", err)
    case errors.Is(err, errors.ErrConflict):
        return NewAPIError(http.StatusConflict, "Resource conflict", err)
    case errors.Is(err, errors.ErrValidation):
        return NewAPIError(http.StatusBadRequest, "Validation failed", err)
    case errors.Is(err, errors.ErrBadRequest):
        return NewAPIError(http.StatusBadRequest, "Bad request", err)
    case errors.Is(err, errors.ErrUnavailable):
        return NewAPIError(http.StatusServiceUnavailable, "Service unavailable", err)
    case errors.Is(err, errors.ErrTimeout):
        return NewAPIError(http.StatusGatewayTimeout, "Request timeout", err)
    default:
        return NewAPIError(http.StatusInternalServerError, "Internal server error", errors.ErrInternal)
    }
}
```

Wrapped errors are detected via `errors.Is()`, so `errors.Wrap(ErrNotFound, "details")` still maps to 404.

### Constructor helpers

```go
api.NotFoundError("User not found")      // 404
api.UnauthorizedError("Invalid token")    // 401
api.ForbiddenError("Access denied")       // 403
api.ConflictError("Email exists")         // 409
api.ValidationError("Invalid email")      // 400
api.BadRequestError("Malformed JSON")     // 400
api.InternalError("msg", underlyingErr)   // 500
api.UnavailableError("DB is down")        // 503
api.TimeoutError("Request too slow")      // 504
```

All support `.WithDetails(map[string]interface{}{...})` for structured error data.

### Handler usage

**Pattern A: Domain-specific sentinel check** (apikeys handler):

```go
resp, err := h.apikeyService.CreateKey(ctx, userID, req)
if err != nil {
    if errors.Is(err, apikeys.ErrMaxKeysExceeded) {
        return &ogen.CreateAPIKeyBadRequest{Code: 400, Message: "Maximum API keys exceeded"}, nil
    }
    if errors.Is(err, apikeys.ErrInvalidScope) {
        return &ogen.CreateAPIKeyBadRequest{Code: 400, Message: err.Error()}, nil
    }
    return &ogen.CreateAPIKeyBadRequest{Code: 500, Message: "Failed to create API key"}, nil
}
```

**Pattern B: Direct ogen response** (most handlers):

```go
if err != nil {
    h.logger.Error("failed to list entities", zap.Error(err))
    return &ogen.Error{Code: 500, Message: "Failed to list entities"}, nil
}
```

**Pattern C: ogen's NewError fallback** (`handler.go`):

```go
func (h *Handler) NewError(ctx context.Context, err error) *ogen.ErrorStatusCode {
    h.logger.Error("Request error", zap.Error(err))
    return &ogen.ErrorStatusCode{
        StatusCode: 500,
        Response:   ogen.Error{Code: 500, Message: "Internal server error"},
    }
}
```

This is the catch-all for unhandled errors from ogen's middleware layer (auth failures, decode errors).

---

## Middleware Errors

Rate limiting (`internal/api/middleware/errors.go`):

```go
func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
    if rateLimitErr, ok := err.(*RateLimitError); ok {
        w.Header().Set("Retry-After", rateLimitErr.RetryAfter.String())
        w.WriteHeader(http.StatusTooManyRequests)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error:   "rate_limit_exceeded",
            Message: "Too many requests",
            Code:    "RATE_LIMIT_EXCEEDED",
        })
        return
    }
    ogenerrors.DefaultErrorHandler(ctx, w, r, err)
}
```

---

## Validation

Input validation in handlers (`internal/validate/`):

```go
limit, err := validate.SafeInt32(params.Limit.Value)
if err != nil {
    return &ogen.SearchForbidden{Code: 400, Message: "Invalid limit"}, nil
}
```

Service-level validation returns simple errors:

```go
if params.Username == "" {
    return nil, fmt.Errorf("username is required")
}
```

---

## Rules

1. **Repositories**: Convert `pgx.ErrNoRows` to domain sentinels. Wrap other errors.
2. **Services**: Use `fmt.Errorf("context: %w", err)`. Define package-level sentinels for domain cases.
3. **Handlers**: Use `errors.Is()` to map domain errors to HTTP status codes. Log with `zap.Error()`.
4. **Never expose**: Internal error details, stack traces, or database error messages in HTTP responses.
5. **Always log**: Log the full error with context before returning a sanitized response.
