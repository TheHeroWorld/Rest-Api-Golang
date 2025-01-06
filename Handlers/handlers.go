package handlers

import (
	"My_Frist_Golang/auth"
	"My_Frist_Golang/db"
	"My_Frist_Golang/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Init_Handlers() {
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := &User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	err := Validation(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		return
	}
	err = db.Registration(&data.Email, &data.Name, &data.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := &AuthUser{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	err := Validation(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		return
	}
	token, err := auth.Auth(data.Email, data.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)

}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = ctx.Value("id")
	data := &task{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	err := Validation(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
	}
	// Переделал под SWITCH, что бы было понятнее какие методы где используются
	switch r.Method {
	case "POST":
		result, err := db.NewTask(id, data.Name, data.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	case "GET":
		taskID := r.URL.Query().Get("Task_id") // Получаем Task_id как строку
		limit := r.URL.Query().Get("Limit")    // Получаем Limit как строку
		data, _ := db.GetAllTasks(id, taskID, limit)
		w.Header().Set("Content-Type", "application/json")
		json_data, _ := json.Marshal(data)
		w.Write(json_data)

	}
}

func ChangeTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id := ctx.Value("id")
	vars := mux.Vars(r)
	// Переделал под SWITCH, что бы было понятнее какие методы где используются
	switch r.Method {
	case "POST":
		data := &task{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(data)
		err := Validation(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error %s", err), http.StatusBadRequest)
		}
		result, err := db.ChangeTask(vars["id"], data.Status, user_id.(float64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	case "DELETE":
		_, err := db.DeleteTask(vars["id"], user_id.(float64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	case "GET":
		data, err := db.GetTask(vars["id"], user_id.(float64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		json_data, _ := json.Marshal(data)
		w.Write(json_data)
	}
}
