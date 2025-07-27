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

	GetAlbumData(database)
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

func GetAlbumData(database *db.DB) {
	rows, err := database.Query("SELECT id, title FROM albums")
	if err != nil {
		fmt.Println("Error querying albums:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Albums:")
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("ID: %d, Title: %s\n", id, title)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error with rows:", err)
	}
}
