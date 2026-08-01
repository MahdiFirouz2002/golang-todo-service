package domain

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput is returned when domain validation fails.
	ErrInvalidInput = errors.New("invalid input")
)
