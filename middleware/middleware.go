package middleware

import (
	"My_Frist_Golang/db"
	"My_Frist_Golang/logging"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = logging.GetLogger() // логгер

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		JWTkey := os.Getenv("KEY")
		key := []byte(JWTkey)                    // это надо убрать в .env
		tokenHeader := r.Header["Authorization"] // ищем хедер Authorization
		if len(tokenHeader) == 0 {
			log.WithFields(logrus.Fields{
				"request_method": r.Method,
				"request_url":    r.URL.String(),
			}).Info("Authorization header is missing")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		tokenString := tokenHeader[0]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Проверяем метод подписания
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				log.WithFields(logrus.Fields{
					"method": token.Header["alg"],
				}).Error("Unexpected signing method")
				return nil, err
			}
			return key, nil
		})
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to parse token")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok { // Достаем из JWT все данные и если все ОК идем дальше
			email := claims["email"]
			id := claims["id"]
			err := db.Findid(id)
			if err != nil {
				log.WithFields(logrus.Fields{
					"user_id": id,
					"email":   email,
					"error":   err.Error(),
				}).Error("Failed to find user by ID")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "email", email) // Суем в контекст r  данные пользавотеля
			ctx = context.WithValue(ctx, "id", id)
			r = r.WithContext(ctx)
			log.WithFields(logrus.Fields{
				"user_id": id,
				"email":   email,
			}).Info("User authenticated successfully")
			next.ServeHTTP(w, r) // Возвращаем к основной функции
		} else {
			log.WithFields(logrus.Fields{
				"request_method": r.Method,
				"request_url":    r.URL.String(),
			}).Warn("Invalid token claims")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		}
	})
}
