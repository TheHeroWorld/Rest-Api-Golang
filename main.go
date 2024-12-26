package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const connStr = "user=postgres password=7458 dbname=test_db sslmode=disable"

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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := &User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	if data.Email == "" || data.Name == "" || data.Password == "" {
		http.Error(w, "Missing fields: email, name or password", http.StatusBadRequest)
	} else {
		err := registration(&data.Email, &data.Name, &data.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := &User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	token, err := auth(data.Email, data.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusUnauthorized)
	}
	w.Header().Set("Authorization", "Bearer "+token)

}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = ctx.Value("id")
	if r.Method == "POST" {
		data := &task{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(data)
		if data.Name == "" || data.Description == "" {
			http.Error(w, "Missing fields: email, name or password", http.StatusBadRequest)
		}
		result, err := NewTask(id, data.Name, data.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Println(result) // НАдо понять как возвращать Idсозданной заявки

	} else {
		data, _ := GetAllTasks(id)
		w.Header().Set("Content-Type", "application/json")
		json_data, _ := json.Marshal(data)
		w.Write(json_data)

	}
}

func ChangeTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if r.Method == "PUT" {
		data := &task{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(data)
		fmt.Println(data)
		_, err := ChangeTusk(vars["id"], data.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else if r.Method == "DELETE" {
		_, err := DeleteTask(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		data, err := GetTask(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(data)
		w.Write(json_data)
	}
}

func main() {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД
	if err != nil {
		panic(err)
	}
	defer db.Close() // отключаемся
	router := mux.NewRouter()
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	protectedRoutes := router.PathPrefix("/tasks").Subrouter()
	protectedRoutes.HandleFunc("", TaskHandler).Methods("POST", "GET")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", ChangeTaskHandler).Methods("PUT", "DELETE", "GET")
	protectedRoutes.Use(authMiddleware)
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
