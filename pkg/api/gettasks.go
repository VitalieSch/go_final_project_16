package api

import (
	"encoding/json"
	"net/http"

	"go1f/pkg/database"
)

// Обработчик для получения списка задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := database.GetTasks()
	if err != nil {
		http.Error(w, `{"error": "ошибка при получении задач"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []database.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]database.Task{"tasks": tasks})

}
