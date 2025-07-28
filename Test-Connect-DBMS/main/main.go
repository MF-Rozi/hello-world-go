package main

import (
	"fmt"
	"os"

	"dev.mfr/db"
	"github.com/joho/godotenv"
)

type Album struct {
	ID     int
	Title  string
	Artist string
	Price  float32
}

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

	GetAlbumsByArtist(database, "John Coltrane")

	GetAlbumById(database, 1)

	GetAlbumByName(database, "Blue Train")

	newAlbum := Album{
		Title:  "Genjirou's First Album",
		Artist: "GenjirouHD",
		Price:  9.99,
	}

	createdAlbum, err := AddAlbum(database, newAlbum)
	if err != nil {
		fmt.Println("Error adding album:", err)
		return
	}
	fmt.Println("Album added successfully:", createdAlbum)

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
	var albums []Album
	rows, err := database.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		fmt.Println("Error querying albums:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Albums:")
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error with rows:", err)
	}
	for _, alb := range albums {
		fmt.Printf("ID: %d, Title: %s, Artist: %s, Price: %.2f\n", alb.ID, alb.Title, alb.Artist, alb.Price)
	}
}

func GetAlbumsByArtist(database *db.DB, artist string) {
	var albums []Album
	rows, err := database.Query("SELECT id, title, artist, price FROM albums WHERE artist = ?", artist)
	if err != nil {
		fmt.Println("Error querying albums by artist:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			fmt.Printf("Error scanning album by artist %q: %v\n", artist, err)
			return
		}
		albums = append(albums, alb)
	}
	fmt.Printf("Albums by %s:\n", artist)
	for _, alb := range albums {
		fmt.Printf("ID: %d, Title: %s, Artist: %s, Price: %.2f\n", alb.ID, alb.Title, alb.Artist, alb.Price)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error with rows:", err)
	}
}

func GetAlbumById(database *db.DB, id int) {
	var album Album
	row := database.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		fmt.Printf("Error querying album by ID %d: %v\n", id, err)
		return
	}
	fmt.Printf("Album by ID %d: Title: %s, Artist: %s, Price: %.2f\n", album.ID, album.Title, album.Artist, album.Price)
}

func GetAlbumByName(database *db.DB, name string) {
	var album Album
	row := database.QueryRow("SELECT id, title, artist, price FROM albums WHERE title = ?", name)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		fmt.Printf("Error querying album by name %s: %v\n", name, err)
		return
	}
	fmt.Printf("Album by Name %s: ID: %d, Artist: %s, Price: %.2f\n", album.Title, album.ID, album.Artist, album.Price)
}

func AddAlbum(database *db.DB, album Album) (Album, error) {
	result, err := database.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return Album{}, fmt.Errorf("error inserting album: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Album{}, fmt.Errorf("error getting last insert ID: %w", err)
	}

	album.ID = int(id)
	return album, nil
}
