package service

import (
	"context"
	"learned/domain"
	"learned/dto"
	"learned/repository"
	"strings"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) userService {
	return userService{repo: repo}
}

func (s userService) CreateUser(ctx context.Context, input dto.CreateUserInput) (*domain.User, error) {
	if input.Name == "" || input.Email == "" {
		return nil, domain.ErrInvalidInput
	}

	user := &domain.User{
		Name:  input.Name,
		Email: input.Email,
	}
	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s userService) PatchUser(
	ctx context.Context,
	id int,
	input dto.PatchUserInput,
) (*domain.User, error) {

	// 1. Validate ID
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// 2. Validate that at least one field is provided
	if input.Name == nil && input.Email == nil {
		return nil, domain.ErrInvalidInput
	}

	// 3. Fetch existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 4. Apply partial updates
	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = strings.ToLower(*input.Email)
	}

	// 5. Persist changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
func (s userService) GetUserByID(
	ctx context.Context,
	id int,
) (*domain.User, error) {

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}
func (s userService) GetAllUsers(
	ctx context.Context,
) ([]*domain.User, error) {
	return s.repo.GetAll(ctx)
}
func (s userService) DeleteUser(
	ctx context.Context,
	id int,
) error {

	if id <= 0 {
		return domain.ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}
func (s userService) UpdateUser(
	ctx context.Context,
	id int,
	input dto.UpdateUserInput,
) (*domain.User, error) {

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	if input.Name == "" || input.Email == "" {
		return nil, domain.ErrInvalidInput
	}

	user := &domain.User{
		ID:    id,
		Name:  input.Name,
		Email: strings.ToLower(input.Email),
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
