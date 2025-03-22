package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"cursor-crash-backend/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=cursoradmin password=Password#1 dbname=cursorcrash port=5432 sslmode='disable'"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Document{})
	err = database.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

	DB = database
}