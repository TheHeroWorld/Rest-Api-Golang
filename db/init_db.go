package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool // переменная пула

func Init_DB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	CONNSTR := os.Getenv("CONNSTR")
	db, err = pgxpool.New(context.Background(), CONNSTR)
	// Открытие пула подключений
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close() // Закрытие пула подключений
	}
}
