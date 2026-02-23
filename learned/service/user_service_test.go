package service

import (
	"context"
	"testing"

	"learned/domain"
	"learned/dto"
)

type mockUserRepository struct {
	CreateFn  func(ctx context.Context, user *domain.User) error
	GetByIDFn func(ctx context.Context, id int) (*domain.User, error)
	UpdateFn  func(ctx context.Context, user *domain.User) error
	DeleteFn  func(ctx context.Context, id int) error
	GetAllFn  func(ctx context.Context) ([]*domain.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	return nil, nil
}
func TestCreateUser_Success(t *testing.T) {
	repo := &mockUserRepository{
		CreateFn: func(ctx context.Context, user *domain.User) error {
			if user.Name != "Sudo" {
				t.Fatalf("unexpected name: %s", user.Name)
			}
			return nil
		},
	}

	svc := NewUserService(repo)

	user, err := svc.CreateUser(context.Background(), dto.CreateUserInput{
		Name:  "Sudo",
		Email: "sudo@test.com",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Name != "Sudo" {
		t.Fatalf("expected name Sudo")
	}
}
func TestCreateUser_InvalidInput(t *testing.T) {
	repo := &mockUserRepository{}

	svc := NewUserService(repo)

	_, err := svc.CreateUser(context.Background(), dto.CreateUserInput{
		Name:  "",
		Email: "",
	})

	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
func TestGetUserByID_NotFound(t *testing.T) {
	repo := &mockUserRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := NewUserService(repo)

	_, err := svc.GetUserByID(context.Background(), 42)

	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
