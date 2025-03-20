package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=cursoradmin password=Password#1 dbname=cursorcrash port=5432 sslmode='disable'"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	database.AutoMigrate(&Document{})
	DB = database
}