package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

const connStr = "user=postgres password=7458 dbname=test_db sslmode=disable"

var ErrInvalidPassword = errors.New("invalid password")
var ErrUserNotFound = errors.New("user not found")

type response struct {
	Tasks []data_task `json:"Tasks"`
}

type data_task struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Deadline_at string `json:"deadline_at"`
}

func Registration(email *string, name *string, password *string) error {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	test := []byte(*password)
	hesh, _ := bcrypt.GenerateFromPassword(test, bcrypt.DefaultCost) // Создаем хэш
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users(email, password, name) values($1, $2, $3)", email, hesh, name) // Запрос к БД
	if err != nil {
		return err
	}
	return nil
}

func FindUser(email string, pass string) (int, error) {
	var storedpassword string
	var id int
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows := db.QueryRow("SELECT password,id FROM users WHERE email = $1", email) // Запрос к БД
	err = rows.Scan(&storedpassword, &id)
	if err != nil {
		return 0, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedpassword), []byte(pass)) // Сравниваем хэши паролей
	if err != nil {
		return 0, ErrInvalidPassword
	}
	return id, nil

}

func NewTask(id any, name string, Description string) ([]data_task, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	time_at := time.Now()
	deadline := time_at.Add(6 * time.Hour)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO tasks(user_id, name, description,created_at, deadline_at) values($1, $2, $3, $4, $5)", id, name, Description, time_at, deadline) // Запрос к БД
	if err != nil {
		return nil, err
	}
	rows := db.QueryRow("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1 ORDER BY tasks DESC LIMIT 1", id) // Запрос к БД
	var p data_task
	err = rows.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	return tasks, nil
}

func GetAllTasks(id any) (response, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1", id) // Запрос к БД
	if err != nil {
		return response{}, err
	}
	var tasks []data_task
	for rows.Next() {
		var p data_task
		err := rows.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
		if err != nil {
			continue
		}
		tasks = append(tasks, p)
	}
	data := response{Tasks: tasks}
	return data, nil
}

func GetTask(id any) ([]data_task, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1", id) // Запрос к БД
	var p data_task
	err = row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	return tasks, nil
}

func DeleteTask(id any) (string, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM tasks WHERE tasks.id = $1", id)
	if err != nil {
		return "", err
	}
	fmt.Println(result)
	return "", nil
}

func ChangeTusk(id any, status string) ([]data_task, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE tasks SET status = $1 where id = $2", status, id) // Запрос к БД
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1", id) // Запрос к БД
	var p data_task
	err = row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	return tasks, nil
}
