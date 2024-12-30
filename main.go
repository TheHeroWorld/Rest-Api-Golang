package main

import (
	"My_Frist_Golang/auth"
	"My_Frist_Golang/db"
	"My_Frist_Golang/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := &User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	if data.Email == "" || data.Name == "" || data.Password == "" {
		http.Error(w, "Missing fields: email, name or password", http.StatusBadRequest)
	} else {
		err := db.Registration(&data.Email, &data.Name, &data.Password)
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
	token, err := auth.Auth(data.Email, data.Password)
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
		result, err := db.NewTask(id, data.Name, data.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	} else {
		data, _ := db.GetAllTasks(id)
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
		result, err := db.ChangeTusk(vars["id"], data.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	} else if r.Method == "DELETE" {
		_, err := db.DeleteTask(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else { // GET
		data, err := db.GetTask(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(data)
		w.Write(json_data)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	protectedRoutes := router.PathPrefix("/tasks").Subrouter() // создаем саб роутер для авторизации
	protectedRoutes.HandleFunc("", TaskHandler).Methods("POST", "GET")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", ChangeTaskHandler).Methods("PUT", "DELETE", "GET")
	protectedRoutes.Use(middleware.AuthMiddleware) // под саб роутер подвязываем мидлвейр авторизации
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
