package database

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Выбор обработчика в зависимости от метода
func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandlerById(w, r)
	case http.MethodPut:
		PutUpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)

	}
}

// Добавление задачи в базу данных
func AddTask(task Task) (int64, error) {
	res, err := database.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Обработчик для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
		return
	}

	if t.Format("20060102") == now.Format("20060102") {
		task.Date = now.Format("20060102")
	} else if t.Before(now) && task.Repeat == "" {
		task.Date = now.Format("20060102")
	} else if t.Before(now) {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = t.Format("20060102")
	}

	if task.Title == "" {
		http.Error(w, `{"error":"не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	part := strings.Split(task.Repeat, " ")

	if part[0] != "y" {
		if part[0] != "d" {
			if part[0] != "" {
				http.Error(w, `{"error":"Неверный формат периодичности задачи"}`, http.StatusBadRequest)
				return
			}
		}
	}

	id, err := AddTask(task)
	if err != nil {
		http.Error(w, `{"error":"ошибка при добавлении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"%d"}`, id)

}

// Функция получения задачи из БД по ID
func GetTaskByID(id string) (Task, error) {
	var t Task
	var err error

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err = database.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, err
	}

	return t, nil
}

// Обработчик для получения задачи методом Get по ID
func GetTaskHandlerById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "ошибка при получении задачи или нет такого id"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

// Функция обновления задачи для БД
func UpdateTask(task *Task) error {
	var err error

	res, err := database.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей,к которым была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

//Функция определения max id в таблице scheduler

func MaxId() (int64, error) {
	var n int64
	row := database.QueryRow(`SELECT max(id) FROM scheduler`)
	err := row.Scan(&n)
	if err != nil {
		panic(err)
	}
	return n, nil
}

// Обработчик обновления задачи методом put
func PutUpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task *Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
		return
	}

	if t.Format("20060102") == now.Format("20060102") {
		task.Date = now.Format("20060102")
	} else if t.Before(now) && task.Repeat == "" {
		task.Date = now.Format("20060102")
	} else if t.Before(now) {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = t.Format("20060102")
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	part := strings.Split(task.Repeat, " ")

	if part[0] != "y" {
		if part[0] != "d" {
			if part[0] != "" {
				http.Error(w, `{"error":"неверный формат периодичности задачи"}`, http.StatusBadRequest)
				return
			}
		}
	}

	if task.ID == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	l, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "не корректный идентификатор"}`, http.StatusBadRequest)
		return
	} else {
		n, err := MaxId()
		if err != nil {
			http.Error(w, `{"error": "не корректный идентификатор"}`, http.StatusBadRequest)
			return
		}
		if l > n {
			http.Error(w, `{"error": "не корректный идентификатор"}`, http.StatusBadRequest)
			return
		}

	}

	UpdateTask(task)

	json.NewEncoder(w).Encode(struct{}{})

}

// Функция удаления задания из БД по ID
func DeleteTaskById(id string) error {

	var err error

	_, err = database.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

// Функция обновления даты задания
func UpdateDate(task *Task) error {

	var err error

	res, err := database.Exec("UPDATE scheduler SET date = ? WHERE id = ?", task.Date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("не корректный id")
	}
	return nil
}

// Обработчик удаления задания
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	_, err := GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"задания по заданному id нет"}`, http.StatusInternalServerError)
		return
	}

	DeleteTaskById(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(struct{}{})

}
