package sse

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/notification"
)

// clientConn represents a connected SSE client.
type clientConn struct {
	id         string
	userID     uuid.UUID
	categories map[notification.EventCategory]bool // nil = all categories
	send       chan []byte
	done       chan struct{}
}

// Broker manages SSE client connections and fans out events.
type Broker struct {
	mu      sync.RWMutex
	clients map[string]*clientConn
	logger  *slog.Logger
}

// NewBroker creates a new SSE broker.
func NewBroker(logger *slog.Logger) *Broker {
	return &Broker{
		clients: make(map[string]*clientConn),
		logger:  logger.With("component", "sse-broker"),
	}
}

// Subscribe registers a new SSE client.
func (b *Broker) Subscribe(userID uuid.UUID, categories []notification.EventCategory) *clientConn {
	conn := &clientConn{
		id:     uuid.Must(uuid.NewV7()).String(),
		userID: userID,
		send:   make(chan []byte, 64),
		done:   make(chan struct{}),
	}

	if len(categories) > 0 {
		conn.categories = make(map[notification.EventCategory]bool, len(categories))
		for _, c := range categories {
			conn.categories[c] = true
		}
	}

	b.mu.Lock()
	b.clients[conn.id] = conn
	b.mu.Unlock()

	b.logger.Info("SSE client connected",
		slog.String("client_id", conn.id),
		slog.String("user_id", userID.String()),
		slog.Int("total_clients", b.ClientCount()),
	)

	return conn
}

// Unsubscribe removes a client connection.
func (b *Broker) Unsubscribe(conn *clientConn) {
	b.mu.Lock()
	if _, ok := b.clients[conn.id]; ok {
		close(conn.done)
		delete(b.clients, conn.id)
	}
	b.mu.Unlock()

	b.logger.Info("SSE client disconnected",
		slog.String("client_id", conn.id),
		slog.Int("total_clients", b.ClientCount()),
	)
}

// Broadcast sends an event to all connected clients that match the category filter.
func (b *Broker) Broadcast(event *notification.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		b.logger.Error("failed to marshal SSE event", slog.Any("error", err))
		return
	}

	category := event.Type.GetCategory()
	msg := formatSSE(event.ID.String(), string(event.Type), data)

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, conn := range b.clients {
		if conn.categories != nil && !conn.categories[category] {
			continue
		}
		select {
		case conn.send <- msg:
		default:
			b.logger.Warn("SSE client buffer full, dropping event",
				slog.String("client_id", conn.id),
				slog.String("event_type", string(event.Type)),
			)
		}
	}
}

// ClientCount returns the number of connected clients.
func (b *Broker) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// formatSSE formats an event as an SSE message.
func formatSSE(id, eventType string, data []byte) []byte {
	msg := "id: " + id + "\nevent: " + eventType + "\ndata: " + string(data) + "\n\n"
	return []byte(msg)
}
