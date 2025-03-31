package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	id := uuid.New()
	tokenString, err := MakeJWT(id, "secret", time.Minute)
	fmt.Printf("token string: %s", tokenString)
	if err != nil {
		t.Errorf("error from MakeJWT")
	}
	validated, err := ValidateJWT(tokenString, "secret")
	if err != nil {
		t.Errorf("error in ValidateJWT, %s", err)
	}
	if validated != id {
		t.Errorf("doesn't match")
	}
}

func TestJWT2(t *testing.T) {
	id := uuid.New()
	tokenString, err := MakeJWT(id, "secret2", time.Second)
	fmt.Printf("token string: %s", tokenString)
	if err != nil {
		t.Errorf("error from MakeJWT")
	}
	time.Sleep(2 * time.Second)
	_, err = ValidateJWT(tokenString, "secret2")
	if err == nil {
		t.Errorf("token should be invalid")
	}
}

func TestJWT3(t *testing.T) {
	id := uuid.New()
	tokenString, err := MakeJWT(id, "secret2", time.Minute)
	fmt.Printf("token string: %s", tokenString)
	if err != nil {
		t.Errorf("error from MakeJWT")
	}
	_, err = ValidateJWT(tokenString, "notsecret")
	if err == nil {
		t.Errorf("token should be invalid")
	}
}
