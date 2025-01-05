package main

import (
	"My_Frist_Golang/db"
	"My_Frist_Golang/handlers"

	_ "github.com/lib/pq"
)

type User struct {
	Email    string `json:"Email"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type task struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	handlers.Init_Handlers()
	err := db.Init_DB()
	if err != nil {
		panic(err)
	}
	defer db.CloseDB()

}
