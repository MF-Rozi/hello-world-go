package main

import (
	"fmt"
	"log"

	"dev.mfr/greetings"
)

func main() {
	fmt.Println("Hello, World!")

	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	message, err := greetings.Greet("Genjirou")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)
}
