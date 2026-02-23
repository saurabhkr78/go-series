package handler

import "learned/domain"

func mapError(err error) (int, string) {
	switch err {
	case domain.ErrInvalidInput:
		return 400, "invalid input"
	case domain.ErrNotFound:
		return 404, "resource not found"
	case domain.ErrConflict:
		return 409, "conflict"
	default:
		return 500, "internal server error"
	}
}
