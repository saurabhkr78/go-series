package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"learned/config"
	"learned/handler"
	"learned/repository"
	"learned/service"
)

func main() {
	// 1. Connect to DB
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer config.CloseDB()

	// 2. Create repository (DB → repository)
	userRepo := repository.NewUserRepository(config.GetDB())

	// 3. Create service (repository → service)
	userService := service.NewUserService(userRepo)

	// 4. Create handler (service → handler)
	userHandler := handler.NewUserHandler(userService)

	// 5. Setup router
	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Patch("/users/{id}", userHandler.PatchUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// 6. Start HTTP server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
