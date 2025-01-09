package db

import (
	"My_Frist_Golang/logging"
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var log = logging.GetLogger() // логгер

var ErrInvalidPassword = errors.New("invalid password")
var ErrUserNotFound = errors.New("user not found")
var ErrNoRows = errors.New("no rows")

type response struct {
	Tasks []data_task `json:"Tasks"`
}

type data_task struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Deadline_at time.Time `json:"deadline_at"`
}

func Registration(email *string, name *string, password *string) error {
	hesh, err := PasswordHesh(password)
	if err != nil {
		log.WithFields(logrus.Fields{
			"email": *email,
			"error": err.Error(),
		}).Error("Password hash creation failed")
		return err
	}
	_, err = db.Exec(context.Background(), "INSERT INTO users(email, password, name) values($1, $2, $3)", email, hesh, name) // Запрос к БД
	if err != nil {
		log.WithFields(logrus.Fields{
			"email": *email,
			"error": err.Error(),
		}).Error("Error inserting user into database")
		return err
	}
	log.WithFields(logrus.Fields{
		"email": *email,
	}).Info("User successfully registered")
	return nil
}

func PasswordHesh(password *string) ([]byte, error) {
	h := []byte(*password)
	hesh, err := bcrypt.GenerateFromPassword(h, bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Password hashing failed")
		return nil, err
	}
	log.Info("Password hashed successfully")
	return hesh, nil
}

func Findid(id any) error {
	rows := db.QueryRow(context.Background(), "SELECT id FROM users WHERE id = $1", id) // Запрос к БД
	err := rows.Scan(&id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("User not found")
		return ErrUserNotFound
	}
	log.WithFields(logrus.Fields{
		"id": id,
	}).Info("User found")
	return nil
}

func FindUser(email string, pass string) (int, error) {
	var storedpassword string
	var id int
	rows := db.QueryRow(context.Background(), "SELECT password,id FROM users WHERE email = $1", email) // Запрос к БД
	err := rows.Scan(&storedpassword, &id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("User not found")
		return 0, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedpassword), []byte(pass)) // Сравниваем хэши паролей
	if err != nil {
		log.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("Invalid password")
		return 0, ErrInvalidPassword
	}
	log.WithFields(logrus.Fields{
		"email": email,
		"id":    id,
	}).Info("User found and password matched")
	return id, nil
}

func NewTask(id any, name string, Description string) ([]data_task, error) {
	time_at := time.Now()
	deadline := time_at.Add(6 * time.Hour)
	_, err := db.Exec(context.Background(), "INSERT INTO tasks(user_id, name, description,created_at, deadline_at) values($1, $2, $3, $4, $5)", id, name, Description, time_at.Format("2006-01-02 15:04:05"), deadline.Format("2006-01-02 15:04:05")) // Запрос к БД
	if err != nil {
		log.WithFields(logrus.Fields{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Error inserting new task into database")
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"user_id": id,
		"name":    name,
	}).Info("New task created successfully")
	rows := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1 ORDER BY tasks DESC LIMIT 1", id) // Запрос к БД
	var p data_task
	err = rows.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		log.WithFields(logrus.Fields{
			"user_id": id,
			"error":   err.Error(),
		}).Error("Error retrieving task after creation")
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	log.Info("Task retrieved successfully after creation")
	return tasks, nil
}

func GetAllTasks(user_id any, task_id string, limit string) (response, error) {
	if task_id == "" {
		task_id = "1" // Если task_id не передан, начинаем с первого ID
	}
	if limit == "" {
		limit = "5" // Если limit не передан, возвращаем 10 записей
	}
	rows, err := db.Query(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks WHERE user_id = $1 AND tasks.id >= $2 ORDER BY tasks.id asc LIMIT $3", user_id, task_id, limit) // КУРСОВАЯ ПАГИНАЦИЯ
	if err != nil {
		log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err.Error(),
		}).Error("Error retrieving tasks")
		return response{}, err
	}
	log.WithFields(logrus.Fields{
		"user_id": user_id,
	}).Info("Tasks retrieved successfully")
	tasks := make([]data_task, 0, 5) // Слайс теперь 5, так как по стандарту лимит 5
	for rows.Next() {
		var p data_task
		err := rows.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error scanning task")
			continue
		}
		tasks = append(tasks, p)
	}
	data := response{Tasks: tasks}
	log.WithFields(logrus.Fields{
		"data": data,
	}).Info("Tasks returned successfully")
	return data, nil
}

func GetTask(id string, user_id float64) ([]data_task, error) {
	row := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1 AND user_id= $2", id, user_id) // Запрос к БД
	var p data_task
	err := row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":      id,
			"user_id": user_id,
			"error":   err.Error(),
		}).Error("Error retrieving task")
		return nil, err
	}
	tasks := make([]data_task, 0, 1) // Создаем слайс с capacitty 1
	tasks = append(tasks, p)
	log.WithFields(logrus.Fields{
		"id":      id,
		"user_id": user_id,
	}).Info("Task retrieved successfully")
	return tasks, nil
}

func DeleteTask(id string, user_id float64) (string, error) { //Подклчается к БД и проверяем ее
	user_id_int := int(user_id)
	result, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE tasks.id = $1 AND tasks.user_id= $2", id, user_id_int)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":      id,
			"user_id": user_id,
			"error":   err.Error(),
		}).Error("Error deleting task")
		return "", err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		log.WithFields(logrus.Fields{
			"id":      id,
			"user_id": user_id,
		}).Warn("No rows affected while deleting task")
		return "", ErrNoRows
	}
	log.WithFields(logrus.Fields{
		"id":      id,
		"user_id": user_id,
	}).Info("Task deleted successfully")
	return "", nil
}

func ChangeTask(id string, status string, user_id float64) ([]data_task, error) {
	_, err := db.Exec(context.Background(), "UPDATE tasks SET status = $1 WHERE id = $2 AND user_id= $3", status, id, user_id) // Запрос к БД
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":      id,
			"user_id": user_id,
			"error":   err.Error(),
		}).Error("Error updating task")
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"id":      id,
		"user_id": user_id,
		"status":  status,
	}).Info("Task status updated successfully")
	row := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1", id) // Запрос к БД
	var p data_task
	err = row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":      id,
			"user_id": user_id,
			"error":   err.Error(),
		}).Error("Error retrieving updated task")
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	log.Info("Task status retrieved successfully after update")
	return tasks, nil
}
