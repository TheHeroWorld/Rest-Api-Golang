package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var db *pgxpool.Pool // переменная пула

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error loading .env file")
		return err
	}
	CONNSTR := os.Getenv("CONNSTR")
	db, err = pgxpool.New(context.Background(), CONNSTR)
	// Открытие пула подключений
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error init database connection pool")
		return err
	}
	log.Info("Database connection pool init successfully")
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close() // Закрытие пула подключений
	}
}
