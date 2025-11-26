package models

import (
	"gorm.io/gorm"
)

// always write ID uint `json:"id" gorm:"primaryKey"` but for know we will use gorm.Model
type Book struct {
	gorm.Model          //automatically includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	Title       string  `json:"title" gorm:"type:varchar(255);not null"`
	Author      string  `json:"author" gorm:"type:varchar(255);not null"`
	Description string  `json:"description" gorm:"type:text"`
	ISBN        string  `json:"isbn" gorm:"uniqueIndex"`
	Price       float64 `json:"price"`
}
