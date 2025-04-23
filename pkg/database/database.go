package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var database *sql.DB

// Инициализация базы данных
func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	database = db

	if install {
		//Создаем таблицу,если ее нет

		schema := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR NOT NULL DEFAULT "",
			comment TEXT DEFAULT "",
			repeat VRCHAR(128) DEFAULT ""
		);
		CREATE INDEX idx_date ON scheduler (date);
		`

		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("не удалось создать таблицу: %v", err)
		}
		fmt.Printf("Файл создан: %s\n", dbFile)
		fmt.Println("Таблица создана")
	} else {
		fmt.Printf("Использован существующий файл: %s\n", dbFile)
	}

	return nil
}
