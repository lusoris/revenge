package api

import "errors"

// API errors.
var (
	// ErrUnauthorized indicates the request is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates the user lacks permission for the action.
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound indicates the resource was not found.
	ErrNotFound = errors.New("not found")
	// ErrModuleDisabled indicates a module is disabled in configuration.
	ErrModuleDisabled = errors.New("module disabled")
)
