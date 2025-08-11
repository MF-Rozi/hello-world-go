package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"dev.mfr/web-service-chi/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var queries *db.Queries
var database *sql.DB

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	database, err = sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("Connected to database successfully!")

	queries = db.New(database)
	// Create table if it doesnt exist
	createTables()

	chi := chi.NewRouter()
	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)

	chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Chi Web Service!")
	})
	chi.Get("/albums", getAlbums)
	chi.Post("/albums", addAlbum)
	chi.Put("/albums/{id}", updateAlbum)
	chi.Get("/albums/name/{name}", findAlbumByName)
	chi.Get("/albums/artist/{artist}", GetAlbumsByArtist)
	chi.Get("/albums/search", getAlbumsByFullTextSearch)
	chi.Delete("/albums/{id}", deleteAlbum)
	chi.Get("/albums/{id}", getAlbumByID)

	http.ListenAndServe(":8080", chi)

}

func getAlbums(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page
	}

	offset := (page - 1) * limit

	var albumsRow []db.GetAlbumsRow
	albumsRow, err = queries.GetAlbums(r.Context(), db.GetAlbumsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching albums: %v", err), http.StatusInternalServerError)
		return
	}
	var Albums []db.Album
	for _, row := range albumsRow {
		Albums = append(Albums, db.Album{
			ID:     row.ID,
			Title:  row.Title,
			Artist: row.Artist,
			Price:  row.Price,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Albums); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding albums: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Fetched albums successfully!")

}
func createTables() {
	// Create the albums table if it doesn't exist
	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS albums (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Error creating albums table: %v\n", err)
	}
	fmt.Println("Albums table created successfully!")
}
func addAlbum(w http.ResponseWriter, r *http.Request) {
	var album db.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding album: %v", err), http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(album.Price, 64)
	if err != nil || album.Title == "" || album.Artist == "" || price <= 0 {
		http.Error(w, "Invalid album data", http.StatusBadRequest)
		return
	}

	newAlbum, err := queries.CreateAlbum(r.Context(), db.CreateAlbumParams{
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating album: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newAlbum); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding new album: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Album added successfully!")
}
func updateAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}
	var album db.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding album: %v", err), http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(album.Price, 64)
	if err != nil || album.Title == "" || album.Artist == "" || price <= 0 {
		http.Error(w, "Invalid album data", http.StatusBadRequest)
		return
	}

	updatedAlbum, err := queries.UpdateAlbum(r.Context(), db.UpdateAlbumParams{
		ID:     int32(id),
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating album: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedAlbum); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated album: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Album updated successfully!")
}
func findAlbumByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Album name is required", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page
	}
	offset := (page - 1) * limit

	albumsRow, err := queries.GetAlbumByTitle(r.Context(), db.GetAlbumByTitleParams{
		Title:  sql.NullString{String: name, Valid: true},
		Limit:  (int32)(limit),
		Offset: (int32)(offset),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching album by name: %v", err), http.StatusInternalServerError)
		return
	}

	if len(albumsRow) == 0 {
		http.Error(w, "Album not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albumsRow); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding album: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Fetched album by name successfully!")
}

func GetAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	artist := chi.URLParam(r, "artist")
	if artist == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page
	}
	offset := (page - 1) * limit

	AlbumRows, err := queries.GetAlbumsByArtist(r.Context(), db.GetAlbumsByArtistParams{
		Artist: sql.NullString{String: artist, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying albums by artist: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(AlbumRows); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding albums by artist: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Fetched albums by %s successfully!\n", artist)
}
func getAlbumsByFullTextSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Default page
	}
	offset := (page - 1) * limit

	albumsRow, err := queries.GetAlbumsByFullTextSearch(r.Context(), db.GetAlbumsByFullTextSearchParams{
		PlaintoTsquery: searchTerm,
		Limit:          int32(limit),
		Offset:         int32(offset),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching albums by full text search: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albumsRow); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding albums by full text search: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Fetched albums by full text search '%s' successfully!\n", searchTerm)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Album ID is required", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid album ID: %v", err), http.StatusBadRequest)
		return
	}
	deletedAlbum, err := queries.GetAlbumByID(r.Context(), int32(idInt))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching deleted album: %v", err), http.StatusInternalServerError)
		return
	}

	if err := queries.DeleteAlbum(r.Context(), (int32)(idInt)); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting album: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Album deleted successfully",
		"album":   deletedAlbum,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	fmt.Println("Album deleted successfully!")
}

func getAlbumByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	album, err := queries.GetAlbumByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching album by ID: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(album); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding album: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Fetched album by ID %d successfully!\n", id)
}
