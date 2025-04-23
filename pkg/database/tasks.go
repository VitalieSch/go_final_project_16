package database

import (
	"errors"
)

const TaskLimit = 50

// Получение списка задач
func GetTasks() ([]Task, error) {

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"

	rows, err := database.Query(query, TaskLimit)
	if err != nil {
		return nil, errors.New("неудалось получить данные")
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("неудалось получить данные")
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
