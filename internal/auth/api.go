package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("no authorization header")
	}
	key, ok := strings.CutPrefix(value, "ApiKey ")
	if !ok {
		return "", errors.New("incorrect format")
	}
	return key, nil
}
