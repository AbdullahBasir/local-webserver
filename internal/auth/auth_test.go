package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf(`HashPassword(%q) = %q returned error: %v`, password, hashed, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "secret"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`HashPassword(%q) returned error: %v`, password, err)
	}

	match, err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Errorf(`CheckPasswordHash(%q, %q) returned error: %v`, password, hashed, err)
	}
	if !match {
		t.Errorf(`CheckPasswordHash(%q, %q) = false, want true`, password, hashed)
	}
}
