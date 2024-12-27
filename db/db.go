package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const connStr = "user=postgres password=7458 dbname=test_db sslmode=disable"

var ErrInvalidPassword = errors.New("invalid password")
var ErrUserNotFound = errors.New("user not found")

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
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users(email, password, name) values($1, $2, $3)", email, password, name)
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
	rows := db.QueryRow("SELECT password,id FROM users WHERE email = $1", email)
	err = rows.Scan(&storedpassword, &id)
	if err != nil {
		return 0, ErrUserNotFound
	}
	if string(pass) == string(storedpassword) {
		return id, nil
	}
	return 0, ErrInvalidPassword
}

func NewTask(id any, name string, Description string) (string, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	time_at := time.Now()
	deadline := time_at.Add(6 * time.Hour)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println(Description)
	_, err = db.Exec("INSERT INTO tasks(user_id, name, description,created_at, deadline_at) values($1, $2, $3, $4, $5)", id, name, Description, time_at, deadline)
	if err != nil {
		return "", err
	}
	return "", nil
}

func GetAllTasks(id any) ([]data_task, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1", id)
	if err != nil {
		return []data_task{}, err
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
	return tasks, nil
}

func GetTask(id any) ([]data_task, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1", id)
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

func ChangeTusk(id any, status string) (string, error) {
	db, err := sql.Open("postgres", connStr) //Подклчается к БД и проверяем ее
	if err != nil {
		panic(err)
	}
	defer db.Close()
	result, err := db.Exec("UPDATE tasks SET status = $1 where id = $2", status, id)
	if err != nil {
		return "", err
	}
	fmt.Println(result.RowsAffected())
	return "", nil
}
