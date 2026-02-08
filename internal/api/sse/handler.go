package sse

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/notification"
)

// Handler serves SSE connections at GET /api/v1/events.
type Handler struct {
	broker       *Broker
	tokenManager auth.TokenManager
	logger       *slog.Logger
}

// NewHandler creates a new SSE handler.
func NewHandler(broker *Broker, tokenManager auth.TokenManager, logger *slog.Logger) *Handler {
	return &Handler{
		broker:       broker,
		tokenManager: tokenManager,
		logger:       logger.With("component", "sse-handler"),
	}
}

// ServeHTTP handles SSE connections.
// Auth: Bearer token in Authorization header or ?token= query param.
// Filtering: ?categories=library,content,system (comma-separated).
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate: try Authorization header first, then query param
	token := ""
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := h.tokenManager.ValidateAccessToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check SSE support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Parse category filters
	var categories []notification.EventCategory
	if cats := r.URL.Query().Get("categories"); cats != "" {
		for _, c := range strings.Split(cats, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				categories = append(categories, notification.EventCategory(c))
			}
		}
	}

	// Disable write timeout for this long-lived connection
	rc := http.NewResponseController(w)
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		h.logger.Warn("failed to disable write deadline", slog.Any("error", err))
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // nginx
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	// Subscribe to events
	conn := h.broker.Subscribe(claims.UserID, categories)
	defer h.broker.Unsubscribe(conn)

	h.logger.Info("SSE connection established",
		slog.String("user_id", claims.UserID.String()),
		slog.String("username", claims.Username),
	)

	// Send keepalive comment every 30s to prevent proxy timeouts
	keepalive := time.NewTicker(30 * time.Second)
	defer keepalive.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-conn.done:
			return
		case msg := <-conn.send:
			if _, err := w.Write(msg); err != nil {
				return
			}
			flusher.Flush()
		case <-keepalive.C:
			// SSE comment (colon prefix) as keepalive
			if _, err := w.Write([]byte(": keepalive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
