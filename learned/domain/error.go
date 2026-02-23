/*
error.go in domain package
Why

# Domain is the center of clean architecture

# All layers are allowed to reference domain errors

Domain errors are stable and long-lived
*/
package domain

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
)
