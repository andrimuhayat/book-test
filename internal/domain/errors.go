package domain

import "errors"

// Sentinel errors for domain-level error handling.
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidData  = errors.New("invalid data")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
)

