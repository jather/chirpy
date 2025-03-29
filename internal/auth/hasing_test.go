package auth

import "testing"

func TestHash(t *testing.T) {
	hash, err := HashPassword("test")
	if err != nil {
		t.Errorf("error while hashing %s", err)
	}
	if CheckPasswordHash(hash, "test") != nil {
		t.Errorf("error while matching, %s", err)
	}
}
