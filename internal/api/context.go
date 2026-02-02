package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Context keys for user authentication data
type contextKey int

const (
	contextKeyUserID contextKey = iota
	contextKeyUsername
	sessionIDKey
)

var (
	// ErrNoUserInContext is returned when no user data is present in context
	ErrNoUserInContext = errors.New("no user in context")
)

// WithUserID stores the user ID in the context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// WithUsername stores the username in the context
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, contextKeyUsername, username)
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(contextKeyUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrNoUserInContext
	}
	return userID, nil
}

// GetUsername retrieves the username from the context
func GetUsername(ctx context.Context) (string, error) {
	username, ok := ctx.Value(contextKeyUsername).(string)
	if !ok {
		return "", ErrNoUserInContext
	}
	return username, nil
}
