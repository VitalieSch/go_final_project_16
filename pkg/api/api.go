package api

import (
	"net/http"

	"go1f/pkg/database"
)

func Init() {

	http.HandleFunc("/api/nextdate", database.NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", GetTasksHandler)
	http.HandleFunc("/api/task/done", DoneTaskHandler)
}
