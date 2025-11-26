package handlers

//body parser populates your struct based on JSON keys and Uses reflection → sets exported fields

// If user sends only "title":Author becomes empty string ("") inside body.We must not overwrite with empty values
//this is called safe partial updates

import (
	"book-management/internal/database"
	"book-management/internal/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// get books
func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	if err := database.DB.Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve books"})
	}
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"}) // first is SELECT * FROM books WHERE id = ? LIMIT 1;
	}
	return c.JSON(book)
}
func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title, author and ISBN are required"})
	}
	if err := database.DB.Create(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create book"})
	}
	return c.Status(fiber.StatusCreated).JSON(book)
}
func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	// Step 1: Fetch existing book
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
	}

	// Step 2: Parse request body
	var body models.Book
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Step 3: Update only changed fields
	if body.Title != "" {
		book.Title = body.Title
	}
	if body.Author != "" {
		book.Author = body.Author
	}
	if body.Description != "" {
		book.Description = body.Description
	}
	if body.Price != 0 {
		book.Price = body.Price
	}

	// Step 4: Save updated book
	if err := database.DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update book"})
	}

	// Step 5: Return updated record
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	//convert id to uint
	if _, err := strconv.Atoi(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid book id"})
	}
	if err := database.DB.Delete(&models.Book{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete book"})
	}

	return c.JSON(fiber.Map{"message": "book deleted successfully"})

}
