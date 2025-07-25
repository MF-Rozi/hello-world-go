package greetings

import (
	"errors"
	"fmt"
)

func Greet(name string) (string, error) {
	if name == "" {
		return "", errors.New("Name cannot be empty")
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}
