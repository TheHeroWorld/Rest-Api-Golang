package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key []byte
	t   *jwt.Token
	s   string
)

func auth(email string, password string) (string, error) {
	id, err := FindUser(email, password)

	if err != nil || id == 0 {
		return "", fmt.Errorf("invalid login or password")
	}
	key = []byte("K0IxiQZOBwHGejUGCTwEz7J9EKi6l1evwEdET/Zy6mg=")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
	})
	s, _ = t.SignedString(key)
	return s, nil
}
