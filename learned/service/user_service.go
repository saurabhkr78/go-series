package service

import (
	"context"
	"errors"
	"learned/domain"
	"learned/dto"
	"learned/repository"
	"strings"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

var (
	ErrInvalidName   = errors.New("name cannot be empty")
	ErrInvalidEmail  = errors.New("email cannot be empty")
	ErrInvalidUserID = errors.New("invalid user id")
	ErrMissingFields = errors.New("all fields are required")
)

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) (*domain.User, error) {
	if input.Name == "" {
		return nil, ErrInvalidName
	}
	if input.Email == "" {
		return nil, ErrInvalidEmail
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
func (s *UserService) PatchUser(
	ctx context.Context,
	id int,
	input dto.PatchUserInput,
) (*domain.User, error) {

	// 1. Validate ID
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	// 2. Validate that at least one field is provided
	if input.Name == nil && input.Email == nil {
		return nil, errors.New("no fields provided for update")
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
func (s *UserService) GetUserByID(
	ctx context.Context,
	id int,
) (*domain.User, error) {

	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	return s.repo.GetByID(ctx, id)
}
func (s *UserService) GetAllUsers(
	ctx context.Context,
) ([]*domain.User, error) {
	return s.repo.GetAll(ctx)
}
func (s *UserService) DeleteUser(
	ctx context.Context,
	id int,
) error {

	if id <= 0 {
		return ErrInvalidUserID
	}

	return s.repo.Delete(ctx, id)
}
func (s *UserService) UpdateUser(
	ctx context.Context,
	id int,
	input dto.UpdateUserInput,
) (*domain.User, error) {

	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	if input.Name == "" || input.Email == "" {
		return nil, ErrMissingFields
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
