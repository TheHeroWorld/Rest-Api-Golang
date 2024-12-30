package auth

import (
	"My_Frist_Golang/db"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key []byte
	t   *jwt.Token
	s   string
)

func Auth(email string, password string) (string, error) {
	key := os.Getenv("KEY")
	id, err := db.FindUser(email, password) // Ищем юзера в базе данных

	if err != nil || id == 0 {
		return "", fmt.Errorf("invalid login or password")
	}
	byte_key := []byte(key)                                      // Надо спрятать в .env
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // Создаем JWT
		"id":    id,
		"email": email,
	})
	s, _ = t.SignedString(byte_key) // Возвращаем JWT
	return s, nil
}
