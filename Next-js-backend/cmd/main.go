package main

import (
	"fmt"

	"dev.mfr/next-js-backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// load environment variables from .env file
	if err := godotenv.Load("internal/config/.env"); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	fmt.Println("Hello, World!")

	server.StartServer()
}
