package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		JWTkey := os.Getenv("KEY")
		fmt.Println(JWTkey)
		key := []byte(JWTkey)                    // это надо убрать в .env
		tokenHeader := r.Header["Authorization"] // ищем хедер Authorization
		if len(tokenHeader) == 0 {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		tokenString := tokenHeader[0]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Проверяем метод подписания
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				return nil, err
			}
			return key, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok { // Достаем из JWT все данные и если все ОК идем дальше
			email := claims["email"]
			id := claims["id"]
			ctx := context.WithValue(r.Context(), "email", email) // Суем в контекст r  данные пользавотеля
			ctx = context.WithValue(ctx, "id", id)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r) // Возвращаем к основной функции

		} else {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		}
	})
}
