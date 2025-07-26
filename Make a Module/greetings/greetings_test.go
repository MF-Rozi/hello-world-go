package greetings

import (
	"regexp"
	"testing"
)

func TestGreet(t *testing.T) {
	name := "Genjirou"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Greet(name)
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Greet("Genjirou") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}
func TestGreetEmpty(t *testing.T) {
	msg, err := Greet("")
	if msg != "" || err == nil {
		t.Fatalf(`Greet("") = %q, %v, want "", error`, msg, err)
	}
}

func TestGreets(t *testing.T) {
	names := []string{"Genjirou", "Hiroshi", "Yuki"}
	messages, err := Greets(names)
	if err != nil {
		t.Fatalf(`Greets(%q) returned error: %v`, names, err)
	}

	for _, name := range names {
		if msg, ok := messages[name]; !ok || msg == "" {
			t.Errorf(`Greets(%q) = %v, want non-empty message for %q`, names, messages, name)
		}
	}
}
