package handlers

import (
	"My_Frist_Golang/auth"
	"My_Frist_Golang/db"
	"My_Frist_Golang/logging"
	"My_Frist_Golang/middleware"
	"My_Frist_Golang/monitoring"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logging.GetLogger() // логгер

func InitHandlers() {
	go monitoring.Monitor() // Запускаем мониторинг в отдельной горутине
	router := mux.NewRouter()

	// Подключаем ErrorHandler middleware
	router.Use(middleware.ErrorHandler)

	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	protectedRoutes := router.PathPrefix("/tasks").Subrouter() // создаем саб роутер для авторизации
	protectedRoutes.HandleFunc("", TaskHandler).Methods("POST", "GET")
	protectedRoutes.HandleFunc("/{id:[0-9]+}", ChangeTaskHandler).Methods("PUT", "DELETE", "GET")
	protectedRoutes.Use(middleware.AuthMiddleware) // под саб роутер подвязываем мидлвейр авторизации
	router.Use(middleware.MonitorMiddleware)       // Мидлверй мониторинга, считает время обрабокти запроса и его статус

	http.Handle("/", router)
	fmt.Println("Server is listening...")
	log.Info("Server is listening on port 8181")
	http.ListenAndServe(":8181", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := DecodeData(&User{}, w, r).(*User)
	log.WithFields(logrus.Fields{
		"email": data.Email,
		"name":  data.Name,
	}).Info("Registration request received")

	err := db.Registration(&data.Email, &data.Name, &data.Password)
	if err != nil {
		// Возвращаем ошибку с помощью нового механизма
		panic(middleware.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Error %s", err)))
	}

	log.WithFields(logrus.Fields{
		"email": data.Email,
	}).Info("Registration successful")
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := DecodeData(&AuthUser{}, w, r).(*AuthUser)
	log.WithFields(logrus.Fields{
		"email": data.Email,
	}).Info("Login request received")
	token, err := auth.Auth(data.Email, data.Password)
	if err != nil {
		log.WithFields(logrus.Fields{
			"email": data.Email,
			"error": err.Error(),
		}).Error("Login failed")
		panic(middleware.NewErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Error %s", err)))
	}
	log.WithFields(logrus.Fields{
		"email": data.Email,
	}).Info("Login successful")
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
		// Возвращаем ошибку с помощью нового механизма
		panic(middleware.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Error %s", err)))
	}

	// Переделал под SWITCH, что бы было понятнее какие методы где используются
	switch r.Method {
	case "POST":
		if data.Name == "" || data.Description == "" {
			// Логируем ошибку
			log.WithFields(logrus.Fields{
				"data": data,
			}).Error("Task creation failed: missing name or description")

			// Отправляем ошибку обратно клиенту
			http.Error(w, "Missing name or description", http.StatusBadRequest)
			return
		}

		result, err := db.NewTask(id, data.Name, data.Description)
		if err != nil {
			log.WithFields(logrus.Fields{ // логи
				"data":  data,
				"error": err.Error(),
			}).Error("Creating task failed")

			// Отправляем ошибку обратно клиенту
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.WithFields(logrus.Fields{ // логи
			"task": data.Name,
		}).Info("Task created successfully")

		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	case "GET":
		taskID := r.URL.Query().Get("Task_id") // Получаем Task_id как строку
		limit := r.URL.Query().Get("Limit")    // Получаем Limit как строку
		data, err := db.GetAllTasks(id, taskID, limit)
		if err != nil {
			log.WithFields(logrus.Fields{ // логи
				"task_id": taskID,
				"limit":   limit,
				"error":   err.Error(),
			}).Error("Get tasks failed")
			panic(middleware.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Error %s", err)))
		}
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
	case "PUT":
		data := DecodeData(&task{}, w, r).(*task)
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
			"status":  data.Status,
			"user_id": user_id,
		}).Info("Update task request received")
		result, err := db.ChangeTask(vars["id"], data.Status, user_id.(float64))
		if err != nil {
			log.WithFields(logrus.Fields{
				"task_id": vars["id"],
				"error":   err.Error(),
			}).Error("Task update failed")
			panic(middleware.NewErrorResponse(http.StatusBadRequest, err.Error()))
		}
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
		}).Info("Task updated successfully")
		json_data, _ := json.Marshal(result)
		w.Write(json_data)

	case "DELETE":
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
			"user_id": user_id,
		}).Info("Delete task request received")
		_, err := db.DeleteTask(vars["id"], user_id.(float64))
		if err != nil {
			log.WithFields(logrus.Fields{
				"task_id": vars["id"],
				"error":   err.Error(),
			}).Error("Task deletion failed")
			panic(middleware.NewErrorResponse(http.StatusBadRequest, err.Error()))
		}
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
		}).Info("Task deleted successfully")
	case "GET":
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
			"user_id": user_id,
		}).Info("Get task request received")
		data, err := db.GetTask(vars["id"], user_id.(float64))
		if err != nil {
			log.WithFields(logrus.Fields{
				"task_id": vars["id"],
				"error":   err.Error(),
			}).Error("Getting task failed")
			panic(middleware.NewErrorResponse(http.StatusBadRequest, err.Error()))
		}
		log.WithFields(logrus.Fields{
			"task_id": vars["id"],
		}).Info("Task retrieved successfully")
		json_data, _ := json.Marshal(data)
		w.Write(json_data)
	}
}
