package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var database *sql.DB

// Инициализация базы данных
func Init(dbFile string) error {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	database = db

	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='scheduler'").Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("не удалось проверить существует ли таблица: %v", err)
	}
	//Создаем таблицу,если ее нет
	if tableExists == 0 {
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
