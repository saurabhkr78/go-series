package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"learned/config"
	"learned/handler"
	"learned/middleware"
	"learned/repository"
	"learned/service"
)

func main() {
	// 1. Connect DB
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer config.CloseDB()

	// 2. Wire dependencies
	userRepo := repository.NewUserRepository(config.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 3. Router
	r := chi.NewRouter()

	// GLOBAL MIDDLEWARE (order matters)
	r.Use(middleware.Recover)
	r.Use(middleware.Logger)

	// USER ROUTES
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Patch("/users/{id}", userHandler.PatchUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// 4. HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 5. Start server
	go func() {
		log.Println("server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server exited properly")
}
