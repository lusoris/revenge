// Package domain contains core business entities and repository interfaces.
package domain

import "errors"

// Common domain errors.
// These are sentinel errors that can be checked with errors.Is.
var (
	// ErrNotFound is returned when the requested entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrUserNotFound is returned when the requested user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrSessionNotFound is returned when the requested session does not exist.
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when the session has expired.
	ErrSessionExpired = errors.New("session expired")

	// ErrOIDCProviderNotFound is returned when the OIDC provider does not exist.
	ErrOIDCProviderNotFound = errors.New("OIDC provider not found")

	// ErrOIDCUserLinkNotFound is returned when the OIDC user link does not exist.
	ErrOIDCUserLinkNotFound = errors.New("OIDC user link not found")

	// ErrDuplicateUsername is returned when attempting to create a user with an existing username.
	ErrDuplicateUsername = errors.New("username already exists")

	// ErrDuplicateEmail is returned when attempting to create a user with an existing email.
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrDuplicateOIDCLink is returned when the OIDC identity is already linked.
	ErrDuplicateOIDCLink = errors.New("OIDC identity already linked")

	// ErrInvalidCredentials is returned when authentication fails.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserDisabled is returned when a disabled user attempts to authenticate.
	ErrUserDisabled = errors.New("user account is disabled")

	// ErrUnauthorized is returned when the user lacks permission.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when access is explicitly denied.
	ErrForbidden = errors.New("forbidden")
)
