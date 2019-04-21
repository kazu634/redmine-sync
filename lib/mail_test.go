package util

import (
	"os"
	"testing"
	"time"
)

var email *Eml

func TestMain(m *testing.M) {
	email = NewEml("./test.eml")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSubject(t *testing.T) {
	if email.Subject() != "Example Message" {
		t.Fatal("Subject should be: Example Message")
	}
}

func TestFrom(t *testing.T) {
	if email.From() != "James Hillyerd <james@inbucket.org>" {
		t.Fatal("From should be: James Hillyerd <james@inbucket.org>")
	}
}

func TestCc(t *testing.T) {
	if email.Cc() != "" {
		t.Fatal("Cc should be empty")
	}
}

func TestDate(t *testing.T) {
	expected := time.Date(2016, 12, 5, 2, 38, 25, 0, time.UTC).String()
	if email.Date().UTC().String() != expected {
		t.Fatalf("Expected %s, but got %s", expected, email.Date().UTC())
	}
}

func TestBody(t *testing.T) {
	if email.Body() != "Text section." {
		t.Fatal("Body should be: Text section.")
	}
}
