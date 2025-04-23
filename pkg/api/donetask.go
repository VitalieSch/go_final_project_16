package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go1f/pkg/database"
)

// Обработчик выполненной задачи, для POST запроса
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task database.Task
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)

	if err != nil {
		http.Error(w, `{"error":"задания по заданному id нет"}`, http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		// Удаление задачи, если она не повторяющаяся
		err = database.DeleteTaskById(task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}

	} else {
		// Обновление даты для повторяющейся задачи
		newDate, err := database.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении даты"}`, http.StatusInternalServerError)
			return
		}

		task.Date = newDate
		database.UpdateDate(&task)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(struct{}{})

}
