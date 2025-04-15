package database

import (
	"encoding/json"
	"net/http"
)

// Получение списка задач
func GetTasks() ([]Task, error) {

	rows, err := database.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Обработчик для получения списка задач

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := GetTasks()
	if err != nil {
		http.Error(w, `{"error": "ошибка при получении задач"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Task{"tasks": tasks})

}
