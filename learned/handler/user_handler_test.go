package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"learned/domain"
	"learned/dto"
)

type mockUserService struct {
	CreateUserFn  func(context.Context, dto.CreateUserInput) (*domain.User, error)
	GetUserByIDFn func(context.Context, int) (*domain.User, error)
	PatchUserFn   func(context.Context, int, dto.PatchUserInput) (*domain.User, error)
	DeleteUserFn  func(context.Context, int) error
}

func (m *mockUserService) CreateUser(ctx context.Context, in dto.CreateUserInput) (*domain.User, error) {
	return m.CreateUserFn(ctx, in)
}
func (m *mockUserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return m.GetUserByIDFn(ctx, id)
}
func (m *mockUserService) PatchUser(ctx context.Context, id int, in dto.PatchUserInput) (*domain.User, error) {
	return m.PatchUserFn(ctx, id, in)
}
func (m *mockUserService) DeleteUser(ctx context.Context, id int) error {
	return m.DeleteUserFn(ctx, id)
}
func TestCreateUser_Success(t *testing.T) {
	mockSvc := &mockUserService{
		CreateUserFn: func(ctx context.Context, in dto.CreateUserInput) (*domain.User, error) {
			return &domain.User{
				ID:    1,
				Name:  "Sudo",
				Email: "sudo@test.com",
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		strings.NewReader(`{"name":"Sudo","email":"sudo@test.com"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestCreateUser_InvalidInput(t *testing.T) {
	mockSvc := &mockUserService{
		CreateUserFn: func(ctx context.Context, in dto.CreateUserInput) (*domain.User, error) {
			return nil, domain.ErrInvalidInput
		},
	}

	handler := NewUserHandler(mockSvc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		strings.NewReader(`{"name":"","email":""}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		GetUserByIDFn: func(ctx context.Context, id int) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	handler := NewUserHandler(mockSvc)

	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserByID)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
