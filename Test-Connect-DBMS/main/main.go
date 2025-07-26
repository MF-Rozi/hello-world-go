package main

import (
	"fmt"
	"os"

	"dev.mfr/db"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Get database connection details from environment variables
	dbConfig := db.Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// Create a new database connection pool
	database, err := db.New(dbConfig)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer database.Close()

	fmt.Println("Connected to database successfully!")

	database.Ping() // Ensure the connection is alive

	GetDatabaseTable(database)

}

func GetDatabaseTable(database *db.DB) {
	tables, err := database.GetTables()
	if err != nil {
		fmt.Println("Error getting tables:", err)
		return
	}

	fmt.Println("Available Tables:")
	for _, table := range tables {
		fmt.Println("-", table)
	}
}
