package routes

import (
	"book-management/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func BookRoutes(app *fiber.App) {
	//group all the routes related to books under book
	book := app.Group("/books")

	book.Get("/", handlers.GetBooks)
	book.Get("/:id", handlers.GetBook)
	book.Post("/", handlers.CreateBook)
	book.Put("/:id", handlers.UpdateBook)
	book.Delete("/:id", handlers.DeleteBook)
}

//now we need to register these routes in the main.go file
