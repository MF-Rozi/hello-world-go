package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

func Greet(name string) (string, error) {
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}

func RandomGreet() string {
	format := []string{
		"Hello, Sir!",
		"Greetings, Friend!",
		"Salutations, Traveler!",
		"Welcome, Guest!",
		"Hi there, Buddy!",
	}
	return format[rand.Intn(len(format))]
}
