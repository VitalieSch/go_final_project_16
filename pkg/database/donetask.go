package database

import (
	"encoding/json"
	"net/http"
	"time"
)

// Функция удаления задачи
func DeleteTask(id string) error {

	_, err := database.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

// Обработчик выполненной задачи, для POST запроса
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task Task
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(id)

	if err != nil {
		http.Error(w, `{"error":"задания по заданному id нет"}`, http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		// Удаление задачи, если она не повторяющаяся
		err = DeleteTask(task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}

	} else {
		// Обновление даты для повторяющейся задачи
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении даты"}`, http.StatusInternalServerError)
			return
		}

		task.Date = newDate
		_, err = database.Exec("UPDATE scheduler SET date = ? WHERE id = ?", task.Date, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(struct{}{})

}
