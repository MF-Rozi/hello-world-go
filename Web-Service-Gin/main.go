package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{ID: 1, Title: "Album One", Artist: "Genjirou", Price: 9.99},
	{ID: 2, Title: "Album Two", Artist: "GenjirouHD", Price: 14.99},
	{ID: 3, Title: "Album Three", Artist: "Genjirou", Price: 19.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8080")

}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
