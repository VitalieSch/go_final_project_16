package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go1f/pkg/database"
)

const DateFmt = "20060102"

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

// Обработчик для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task database.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFmt)
	}

	t, err := time.Parse(DateFmt, task.Date)
	if err != nil {
		http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
		return
	}

	if t.Format("20060102") == now.Format(DateFmt) {
		task.Date = now.Format(DateFmt)
	} else if t.Before(now) && task.Repeat == "" {
		task.Date = now.Format(DateFmt)
	} else if t.Before(now) {
		nextDate, err := database.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else {
		task.Date = t.Format(DateFmt)
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

	id, err := database.AddTask(task)
	if err != nil {
		http.Error(w, `{"error":"ошибка при добавлении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"%d"}`, id)

}

// Обработчик для получения задачи методом Get по ID
func GetTaskHandlerById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "ошибка при получении задачи или нет такого id"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

// Обработчик обновления задачи методом put
func PutUpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task *database.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}
	l, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error": "ошибка обработки данных"}`, http.StatusBadRequest)
		return
	}

	if l > database.LastId() {
		http.Error(w, `{"error": "некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFmt)
	}
	t, err := time.Parse(DateFmt, task.Date)
	if err != nil {
		http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
		return
	}

	if t.Format("20060102") == now.Format(DateFmt) {
		task.Date = now.Format(DateFmt)
	} else if t.Before(now) && task.Repeat == "" {
		task.Date = now.Format(DateFmt)
	} else if t.Before(now) {
		task.Date, err = database.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"неверный формат двты"}`, http.StatusBadRequest)
			return
		}

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

	database.UpdateTask(task)

	json.NewEncoder(w).Encode(struct{}{})

}

// Обработчик удаления задания
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	_, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"задания по заданному id нет"}`, http.StatusInternalServerError)
		return
	}

	database.DeleteTaskById(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(struct{}{})

}
