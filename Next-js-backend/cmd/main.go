package main

import (
	"fmt"

	"dev.mfr/next-js-backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// to start the server go to root directory and run `go run cmd/main.go`
	fmt.Println("Starting the application...")
	env_url := "internal/config/.env"
	fmt.Println("Loading environment variables from: ", env_url)
	// load environment variables from .env file
	if err := godotenv.Load(env_url); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	fmt.Println("Hello, World!")

	server.StartServer()
}
