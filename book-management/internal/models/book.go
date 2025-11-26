package models

import (
	"gorm.io/gorm"
)

// always write ID uint `json:"id" gorm:"primaryKey"` but for know we will use gorm.Model
type Book struct {
	gorm.Model          //automatically includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ISBN        string  `json:"isbn" gorm:"uniqueIndex"`
	Price       float64 `json:"price"`
}
