package static

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {

	// 1️. Create a new Fiber app
	app := fiber.New()

	// 2️. Serve static files (HTML, CSS, JS) from ./static
	app.Static("/", "./static")

	// ===================== FORM HANDLER (POST) =====================
	app.Post("/form", func(c *fiber.Ctx) error {

		// 1️. Parse form data (Fiber parses automatically)
		name := c.FormValue("name")
		address := c.FormValue("address")

		// 2️. Send response
		return c.SendString(
			fmt.Sprintf("POST request successful\nHello %s\nAddress = %s\n", name, address),
		)
	})

	// ===================== HELLO HANDLER (GET) =====================
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello from server!")
	})

	// 3️. Start the server
	fmt.Println("Server is running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
