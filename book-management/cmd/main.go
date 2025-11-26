package main

import (
	"book-management/internal/database"
	"book-management/internal/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	//connect database
	database.ConnectDB()
	//register book routes
	routes.BookRoutes(app)

	app.Listen(":8000")
}
