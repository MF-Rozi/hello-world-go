package main

import (
	"fmt"
	"net/http"
	"os"

	"dev.mfr/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

//	var albums = []Album{
//		{ID: 1, Title: "Album One", Artist: "Genjirou", Price: 9.99},
//		{ID: 2, Title: "Album Two", Artist: "GenjirouHD", Price: 14.99},
//		{ID: 3, Title: "Album Three", Artist: "Genjirou", Price: 19.99},
//	}
var database *db.DB
var err error

func main() {

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
	database, err = db.New(dbConfig)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer database.Close()
	fmt.Println("Connected to database successfully!")

	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8080")

}

func getAlbums(c *gin.Context) {
	var albums []Album
	albm, err := database.Query("SELECT id, title, artist, price FROM albums")

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch albums"})
		return
	}
	for albm.Next() {
		var album Album
		if err := albm.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan album"})
			return
		}
		albums = append(albums, album)
	}
	c.IndentedJSON(http.StatusOK, albums)
}
