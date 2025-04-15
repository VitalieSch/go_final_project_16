package api

import (
	"go1f/pkg/database"
	"net/http"
)

func Init() {

	http.HandleFunc("/api/nextdate", database.NextDateHandler)
	http.HandleFunc("/api/task", database.TaskHandler)
	http.HandleFunc("/api/tasks", database.GetTasksHandler)
	http.HandleFunc("/api/task/done", database.DoneTaskHandler)
}
