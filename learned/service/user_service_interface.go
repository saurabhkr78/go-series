package service

import (
	"context"

	"learned/domain"
	"learned/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, input dto.CreateUserInput) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	PatchUser(ctx context.Context, id int, input dto.PatchUserInput) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}
