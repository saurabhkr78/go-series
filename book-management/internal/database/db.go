package database

import (
	"book-management/internal/config"
	"book-management/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

//1. create a variable to hold the database connection object which will be of pointer type gorm.DB

var DB *gorm.DB

func ConnectDB() {
	//load the config to get the database credentials
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	//connect to the database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to the database", err)
	}

	//run the database migrations
	err = database.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Could not migrate the database", err)
	}

	//assign the database connection object to the DB variable
	DB = database

	fmt.Println("Connected to the database successfully and auto migrated")
}
