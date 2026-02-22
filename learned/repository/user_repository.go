package repository

import (
	"context"

	"learned/domain"
)

// what can be done eith domain enttity i.e user
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
}
