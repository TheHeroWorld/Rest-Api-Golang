package db

import (
	"context"
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

var db *pgxpool.Pool
var err error

func Init_DB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	CONNSTR := os.Getenv("CONNSTR")
	db, err = pgxpool.New(context.Background(), CONNSTR) // Открытие пула подключений
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close() // Закрытие пула подключений
	}
}

func Registration(email *string, name *string, password *string) error {
	hesh, err := PasswordHesh(password)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), "INSERT INTO users(email, password, name) values($1, $2, $3)", email, hesh, name) // Запрос к БД
	if err != nil {
		return err
	}
	return nil
}

func PasswordHesh(password *string) ([]byte, error) {
	h := []byte(*password)
	hesh, err := bcrypt.GenerateFromPassword(h, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hesh, nil
}

func FindUser(email string, pass string) (int, error) {
	var storedpassword string
	var id int
	rows := db.QueryRow(context.Background(), "SELECT password,id FROM users WHERE email = $1", email) // Запрос к БД
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
	time_at := time.Now()
	deadline := time_at.Add(6 * time.Hour)
	_, err = db.Exec(context.Background(), "INSERT INTO tasks(user_id, name, description,created_at, deadline_at) values($1, $2, $3, $4, $5)", id, name, Description, time_at.Format("2006-01-02 15:04:05"), deadline.Format("2006-01-02 15:04:05")) // Запрос к БД
	if err != nil {
		return nil, err
	}
	rows := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1 ORDER BY tasks DESC LIMIT 1", id) // Запрос к БД
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
	rows, err := db.Query(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE users.id = $1", id) // Запрос к БД
	if err != nil {
		return response{}, err
	}
	tasks := make([]data_task, 0, 2) // Пустой слайст с обьемом 2, что бы лишний раз не увеличивать слайс если 2 заявки
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
func GetTask(id string, user_id float64) ([]data_task, error) {
	row := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1 AND user_id= $2", id, user_id) // Запрос к БД
	var p data_task
	err = row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		return nil, err
	}

	tasks := make([]data_task, 0, 1) // Создаем слайс с capacitty 1
	tasks = append(tasks, p)
	return tasks, nil
}

// Надо сделать так что бы удалять данные администратор
func DeleteTask(id string, user_id float64) (string, error) { //Подклчается к БД и проверяем ее
	user_id_int := int(user_id)
	result, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE tasks.id = $1 AND tasks.user_id= $2", id, user_id_int)
	if err != nil {
		return "", err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return "", ErrNoRows
	}
	return "", nil
}

// Надо сделать так что бы менять данные администратор
func ChangeTask(id string, status string, user_id float64) ([]data_task, error) {
	_, err = db.Exec(context.Background(), "UPDATE tasks SET status = $1 WHERE id = $2 AND user_id= $3", status, id, user_id) // Запрос к БД
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(context.Background(), "SELECT tasks.id, status, tasks.name, description, created_at,deadline_at FROM tasks,users WHERE tasks.id = $1", id) // Запрос к БД
	var p data_task
	err = row.Scan(&p.ID, &p.Status, &p.Name, &p.Description, &p.CreatedAt, &p.Deadline_at)
	if err != nil {
		return nil, err
	}
	var tasks []data_task
	tasks = append(tasks, p)
	return tasks, nil
}
