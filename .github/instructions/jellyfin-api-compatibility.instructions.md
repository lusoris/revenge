---
applyTo: "**/internal/api/**/*.go"
---

# Jellyfin API Compatibility Reference

> Translating C# Jellyfin API patterns to Go

## Goal: 100% API Compatibility

Every endpoint must match the original C# Jellyfin:

- **Same route** (path, method, query params)
- **Same response structure** (JSON field names, types, nesting)
- **Same behavior** (validation, defaults, error codes)

## C# Controller Pattern → Go Handler

### Original C# Controller

```csharp
// Jellyfin.Api/Controllers/UserController.cs
[Route("Users")]
public class UserController : BaseJellyfinApiController
{
    private readonly IUserManager _userManager;
    private readonly ISessionManager _sessionManager;

    public UserController(
        IUserManager userManager,
        ISessionManager sessionManager)
    {
        _userManager = userManager;
        _sessionManager = sessionManager;
    }

    [HttpGet("{userId}")]
    [Authorize(Policy = Policies.IgnoreParentalControl)]
    [ProducesResponseType(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public ActionResult<UserDto> GetUserById([FromRoute, Required] Guid userId)
    {
        var user = _userManager.GetUserById(userId);
        if (user is null)
        {
            return NotFound("User not found");
        }

        return _userManager.GetUserDto(user, HttpContext.GetNormalizedRemoteIP().ToString());
    }
}
```

### Go Equivalent

```go
// internal/api/handlers/user.go
package handlers

type UserHandler struct {
    userService    *service.UserService
    sessionService *service.SessionService
}

func NewUserHandler(
    userService *service.UserService,
    sessionService *service.SessionService,
) *UserHandler {
    return &UserHandler{
        userService:    userService,
        sessionService: sessionService,
    }
}

// GET /Users/{userId}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    userID, err := uuid.Parse(r.PathValue("userId"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    user, err := h.userService.GetUserByID(r.Context(), userID)
    if errors.Is(err, service.ErrUserNotFound) {
        writeError(w, http.StatusNotFound, "User not found")
        return
    }
    if err != nil {
        slog.Error("failed to get user", "error", err, "user_id", userID)
        writeError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    dto := h.userService.ToDTO(user, getRemoteIP(r))
    writeJSON(w, http.StatusOK, dto)
}
```

## Route Registration

```go
// internal/api/router.go
func NewRouter(h *handlers.Handlers, mw *middleware.Middleware) http.Handler {
    mux := http.NewServeMux()

    // Users endpoints - match C# [Route("Users")]
    mux.HandleFunc("GET /Users", mw.Auth(h.User.GetUsers))
    mux.HandleFunc("GET /Users/Public", h.User.GetPublicUsers)
    mux.HandleFunc("GET /Users/{userId}", mw.Auth(h.User.GetUserByID))
    mux.HandleFunc("DELETE /Users/{userId}", mw.RequireAdmin(h.User.DeleteUser))
    mux.HandleFunc("POST /Users/New", mw.RequireAdmin(h.User.CreateUser))
    mux.HandleFunc("POST /Users/AuthenticateByName", h.User.AuthenticateByName)
    mux.HandleFunc("POST /Users/Password", mw.Auth(h.User.UpdatePassword))
    mux.HandleFunc("POST /Users/{userId}/Policy", mw.RequireAdmin(h.User.UpdatePolicy))
    mux.HandleFunc("GET /Users/Me", mw.Auth(h.User.GetCurrentUser))

    return mux
}
```

## Common API Patterns

### Parameter Sources

| C# Attribute   | Go Equivalent                          |
| -------------- | -------------------------------------- |
| `[FromRoute]`  | `r.PathValue("param")`                 |
| `[FromQuery]`  | `r.URL.Query().Get("param")`           |
| `[FromBody]`   | `json.NewDecoder(r.Body).Decode(&req)` |
| `[FromHeader]` | `r.Header.Get("X-Header")`             |

### Response Types

| C# Return          | Go Equivalent                |
| ------------------ | ---------------------------- |
| `Ok(result)`       | `writeJSON(w, 200, result)`  |
| `NoContent()`      | `w.WriteHeader(204)`         |
| `NotFound()`       | `writeError(w, 404, "...")`  |
| `BadRequest()`     | `writeError(w, 400, "...")`  |
| `Forbid()`         | `writeError(w, 403, "...")`  |
| `StatusCode(code)` | `writeError(w, code, "...")` |

### Authorization Policies

```csharp
// C# Policy attributes
[Authorize]                                    // Requires authentication
[Authorize(Policy = Policies.RequiresElevation)]  // Admin only
[Authorize(Policy = Policies.IgnoreParentalControl)]
[Authorize(Policy = Policies.FirstTimeSetupOrElevated)]
```

```go
// Go middleware equivalents
mw.Auth(handler)           // Requires authentication
mw.RequireAdmin(handler)   // Admin only
mw.IgnoreParental(handler) // Skip parental controls
mw.FirstTimeOrAdmin(handler)
```

## Key Controllers to Implement

### Priority 1 - Authentication & Users

- `UserController` - User CRUD, auth, passwords
- `ApiKeyController` - API key management
- `QuickConnectController` - Quick connect flow

### Priority 2 - Library & Media

- `ItemsController` - Media items
- `LibraryController` - Library management
- `UserLibraryController` - User-specific library access
- `UserViewsController` - User views/collections

### Priority 3 - Playback

- `MediaInfoController` - Media information
- `PlaystateController` - Play state tracking
- `SessionController` - Session management

### Priority 4 - Metadata

- `ArtistsController`, `GenresController`, `PersonsController`
- `StudiosController`, `YearsController`
- `RemoteImageController`, `ImageController`

### Priority 5 - Streaming

- `AudioController`, `VideosController`
- `DynamicHlsController`, `HlsSegmentController`
- `SubtitleController`

## DTO Mapping

### C# UserDto

```csharp
public class UserDto
{
    public Guid Id { get; set; }
    public string? Name { get; set; }
    public string? ServerId { get; set; }
    public bool HasPassword { get; set; }
    public bool HasConfiguredPassword { get; set; }
    public DateTime? LastLoginDate { get; set; }
    public DateTime? LastActivityDate { get; set; }
    public UserConfiguration? Configuration { get; set; }
    public UserPolicy? Policy { get; set; }
}
```

### Go Equivalent

```go
// internal/domain/user.go
type UserDTO struct {
    ID                    uuid.UUID         `json:"Id"`
    Name                  string            `json:"Name,omitempty"`
    ServerID              string            `json:"ServerId,omitempty"`
    HasPassword           bool              `json:"HasPassword"`
    HasConfiguredPassword bool              `json:"HasConfiguredPassword"`
    LastLoginDate         *time.Time        `json:"LastLoginDate,omitempty"`
    LastActivityDate      *time.Time        `json:"LastActivityDate,omitempty"`
    Configuration         *UserConfiguration `json:"Configuration,omitempty"`
    Policy                *UserPolicy       `json:"Policy,omitempty"`
}
```

**CRITICAL**: JSON field names must match C# exactly (PascalCase)!

## Common Types

### GUID/UUID

```go
import "github.com/google/uuid"

// Parse from string
id, err := uuid.Parse(r.PathValue("userId"))

// Check for empty/nil
if id == uuid.Nil {
    return ErrInvalidID
}
```

### Nullable Types

```csharp
// C# nullable
public Guid? UserId { get; set; }
public DateTime? LastLogin { get; set; }
```

```go
// Go pointers for nullable
type Request struct {
    UserID    *uuid.UUID `json:"UserId,omitempty"`
    LastLogin *time.Time `json:"LastLogin,omitempty"`
}
```

### Date Format

```go
// Jellyfin date format
const JellyfinDateFormat = "2006-01-02T15:04:05.0000000Z"

func formatDate(t time.Time) string {
    return t.UTC().Format(JellyfinDateFormat)
}
```

## Helper Functions

```go
// internal/api/helpers.go

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getRemoteIP(r *http.Request) string {
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        return strings.Split(xff, ",")[0]
    }
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
}

func parseUUID(s string) (uuid.UUID, error) {
    return uuid.Parse(s)
}

func parseOptionalUUID(s string) *uuid.UUID {
    if s == "" {
        return nil
    }
    id, err := uuid.Parse(s)
    if err != nil {
        return nil
    }
    return &id
}

func parseBool(s string) bool {
    return s == "true" || s == "True" || s == "1"
}

func parseInt(s string, def int) int {
    if s == "" {
        return def
    }
    v, err := strconv.Atoi(s)
    if err != nil {
        return def
    }
    return v
}
```

## Query Result Pattern

```csharp
// C# QueryResult
public class QueryResult<T>
{
    public IReadOnlyList<T> Items { get; set; }
    public int TotalRecordCount { get; set; }
    public int StartIndex { get; set; }
}
```

```go
// Go equivalent
type QueryResult[T any] struct {
    Items            []T `json:"Items"`
    TotalRecordCount int `json:"TotalRecordCount"`
    StartIndex       int `json:"StartIndex"`
}

func NewQueryResult[T any](items []T, total, start int) QueryResult[T] {
    if items == nil {
        items = []T{} // Never null, always empty array
    }
    return QueryResult[T]{
        Items:            items,
        TotalRecordCount: total,
        StartIndex:       start,
    }
}
```

## Request Body Pattern

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserByName
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Validate (match C# [Required] attributes)
    if strings.TrimSpace(req.Name) == "" {
        writeError(w, http.StatusBadRequest, "Name is required")
        return
    }

    // Process...
}

type CreateUserByName struct {
    Name     string  `json:"Name"`
    Password *string `json:"Password,omitempty"`
}
```

## Testing API Compatibility

```go
func TestGetUserByID_MatchesJellyfin(t *testing.T) {
    // Response must match Jellyfin exactly
    want := `{
        "Id": "...",
        "Name": "TestUser",
        "HasPassword": true,
        "HasConfiguredPassword": true
    }`

    // Compare JSON structure, not just values
    var wantMap, gotMap map[string]any
    json.Unmarshal([]byte(want), &wantMap)
    json.Unmarshal(rec.Body.Bytes(), &gotMap)

    // All keys must match
    for k := range wantMap {
        if _, ok := gotMap[k]; !ok {
            t.Errorf("missing key: %s", k)
        }
    }
}
```

## Reference Sources

When implementing an endpoint:

1. Find the C# controller in `Jellyfin.Api/Controllers/`
2. Note the route, attributes, parameters, response types
3. Check DTOs in `MediaBrowser.Model/` and `Jellyfin.Api/Models/`
4. Check interfaces in `MediaBrowser.Controller/`
5. Check implementations in `Jellyfin.Server.Implementations/`
6. Run original Jellyfin and capture real API responses for testing
