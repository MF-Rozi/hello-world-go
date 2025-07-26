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

	return RandomGreet(name), nil
	// return fmt.Sprintf("Hello, %s!", name), nil
}

func Greets(names []string) (map[string]string, error) {
	messages := make(map[string]string)
	for _, name := range names {
		message, err := Greet(name)
		if err != nil {
			return nil, err
		}
		messages[name] = message
	}
	return messages, nil
}
func RandomGreet(name string) string {
	format := []string{
		fmt.Sprintf("Hello, %s!", name),
		fmt.Sprintf("Greetings, %s!", name),
		fmt.Sprintf("Salutations, %s!", name),
		fmt.Sprintf("Welcome, %s!", name),
		fmt.Sprintf("Hi there, %s!", name),
	}
	return format[rand.Intn(len(format))]
}
