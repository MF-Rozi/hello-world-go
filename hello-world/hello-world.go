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

	fmt.Println(greetings.RandomGreet())

	names := []string{"Genjirou", "Hiroshi", "Yuki"}

	messages, err := greetings.Greets(names)
	if err != nil {
		log.Fatal(err)
	}
	for name, message := range messages {
		fmt.Printf("%s: %s\n", name, message)
	}
}
