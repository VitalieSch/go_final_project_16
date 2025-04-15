package main

import (
	"database/sql"
	"go1f/pkg/database"
	"go1f/pkg/server"
	"log"
)

func main() {
	var db *sql.DB

	err := database.Init("scheduler.db")
	if err != nil {
		log.Fatalf("не удалось получить БД: %v\n", err)
	}

	defer db.Close()
	server.Run()

}
