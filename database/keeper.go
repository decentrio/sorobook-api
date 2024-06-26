package database

import (
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)


func NewDBHandler() *gorm.DB {
	db := createConnection()
	return db
}

// create connection with postgres db
func createConnection() *gorm.DB {
	sqlUrl, ok := os.LookupEnv("READONLY_URL")

	if !ok {
		log.Fatalf("Error get READONLY_URL")
	}

	// Open the connection
	db, err := gorm.Open(postgres.Open(sqlUrl), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	log.Println("Connected to PSQL:", db)

	return db
}