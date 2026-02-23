package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"

	"learned/config"
	"learned/domain"
	"learned/handler"
	"learned/repository"
	"learned/service"
)

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	// connect real DB (test DB!)
	if err := config.ConnectDB(); err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// clean DB before each test run
	_, err := config.GetDB().Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY")
	if err != nil {
		t.Fatalf("failed to clean db: %v", err)
	}

	// real wiring
	repo := repository.NewUserRepository(config.GetDB())
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	r := chi.NewRouter()
	r.Post("/users", h.CreateUser)
	r.Get("/users/{id}", h.GetUserByID)
	r.Patch("/users/{id}", h.PatchUser)
	r.Delete("/users/{id}", h.DeleteUser)

	return httptest.NewServer(r)
}
func teardownTestServer() {
	config.CloseDB()
}
func TestUserLifecycle(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()
	defer teardownTestServer()

	client := server.Client()

	// ---------- CREATE USER ----------
	createBody := map[string]string{
		"name":  "Sudo",
		"email": "sudo@test.com",
	}

	bodyBytes, _ := json.Marshal(createBody)

	resp, err := client.Post(
		server.URL+"/users",
		"application/json",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created domain.User
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response failed: %v", err)
	}

	if created.ID == 0 {
		t.Fatalf("expected generated user id")
	}

	// ---------- GET USER ----------
	resp, err = client.Get(
		server.URL + "/users/" + strconv.Itoa(created.ID),
	)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var fetched domain.User
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get response failed: %v", err)
	}

	if fetched.Email != "sudo@test.com" {
		t.Fatalf("unexpected email: %s", fetched.Email)
	}

	// ---------- DELETE USER ----------
	req, _ := http.NewRequest(
		http.MethodDelete,
		server.URL+"/users/"+strconv.Itoa(created.ID),
		nil,
	)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	// ---------- GET AFTER DELETE (NOT FOUND) ----------
	resp, err = client.Get(
		server.URL + "/users/" + strconv.Itoa(created.ID),
	)
	if err != nil {
		t.Fatalf("get-after-delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
