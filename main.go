package main

import (
	"My_Frist_Golang/db"
	"My_Frist_Golang/handlers"

	_ "github.com/lib/pq"
)

func main() {
	err := db.Init_DB() // Если поменять местами с Handlers, то база данных не подключается почемУ???
	if err != nil {
		panic(err)
	}
	handlers.Init_Handlers()

	defer db.CloseDB()

}
