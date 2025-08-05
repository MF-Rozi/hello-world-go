package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	router.GET("/albums/:id", getAlbumByID)
	router.GET("/albums/name/:name", GetAlbumByName)
	router.POST("/albums/", AddAlbum)
	router.PUT("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)
	router.GET("/albums/search", FindAlbumByFullTextSearch)

	router.Run("localhost:8080")
	fmt.Println("Server running on http://localhost:8080")

}

func getAlbums(c *gin.Context) {
	// Get pagination parameters from query string
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(400, gin.H{"error": "Invalid limit (1-100)"})
		return
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var total int
	countRow := database.QueryRow("SELECT COUNT(*) FROM albums")
	if err := countRow.Scan(&total); err != nil {
		c.JSON(500, gin.H{"error": "Failed to count albums"})
		return
	}

	// Get paginated albums
	var albums []Album
	albm, err := database.Query("SELECT id, title, artist, price FROM albums LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer albm.Close()

	for albm.Next() {
		var album Album
		if err := albm.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan album"})
			return
		}
		albums = append(albums, album)
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	// Return paginated response
	c.IndentedJSON(http.StatusOK, gin.H{
		"data": albums,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	})
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	var album Album

	row := database.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		c.JSON(404, gin.H{"error": "Album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func GetAlbumByName(c *gin.Context) {
	name := c.Param("name")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	searchTerm := "%" + name + "%"

	// Get total count for this search
	var total int
	countRow := database.QueryRow("SELECT COUNT(*) FROM albums WHERE title LIKE ?", searchTerm)
	countRow.Scan(&total)

	// Get paginated results
	var albums []Album
	rows, err := database.Query("SELECT id, title, artist, price FROM albums WHERE title LIKE ? LIMIT ? OFFSET ?", searchTerm, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		albums = append(albums, album)
	}

	totalPages := (total + limit - 1) / limit

	c.IndentedJSON(http.StatusOK, gin.H{
		"data": albums,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	})
}

func AddAlbum(c *gin.Context) {
	var newAlbum Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(400, gin.H{"error": "Invalid album data"})
		return
	}

	result, err := database.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add album"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve album ID"})
		return
	}
	newAlbum.ID = int(id)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func updateAlbum(c *gin.Context) {
	id := c.Param("id")
	var updatedAlbum Album
	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.JSON(400, gin.H{"error": "Invalid album data"})
		return
	}

	integerid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid album ID"})
		return
	}
	updatedAlbum.ID = integerid

	result, err := database.Exec("UPDATE albums SET title = ?, artist = ?, price = ? WHERE id = ?", updatedAlbum.Title, updatedAlbum.Artist, updatedAlbum.Price, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update album"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedAlbum)
}
func deleteAlbum(c *gin.Context) {
	id := c.Param("id")
	var album Album
	row := database.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		c.JSON(404, gin.H{"error": "Album not found"})
		return
	}

	result, err := database.Exec("DELETE FROM albums WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete album"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album deleted successfully", "album": album})
}
func FindAlbumByFullTextSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	var albums []Album
	rows, err := database.Query("SELECT id, title, artist, price FROM albums WHERE MATCH(title, artist) AGAINST(? IN NATURAL LANGUAGE MODE)", query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search albums"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		albums = append(albums, album)
	}
	c.IndentedJSON(http.StatusOK, albums)
}
