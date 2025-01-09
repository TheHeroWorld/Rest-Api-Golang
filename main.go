package main

import (
	"My_Frist_Golang/db"
	"My_Frist_Golang/handlers"
	"My_Frist_Golang/logging"

	_ "github.com/lib/pq"
)

func main() {
	logging.InitLog()
	err := db.InitDB() // Если поменять местами с Handlers, то база данных не подключается почемУ???
	if err != nil {
		panic(err)
	}
	handlers.InitHandlers()
	defer db.CloseDB()

}
