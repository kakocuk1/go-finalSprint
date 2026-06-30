package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT, 
    date    CHAR(8) NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler (date);
`

var DB *sql.DB /* хранение глобальной переменной для подключения к базе данных.
такой подход считается приемлемым и для небольших production-сервисов? */

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true // если install равен true, после открытия БД требуется выполнить sql-запрос с CREATE TABLE и CREATE INDEX
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := db.Exec(schema); err != nil { // создаем schema, если БД создается впервые, если был то просто открываем без создания schema
			db.Close() // правки после ревью - закрываем коннект перед выходом
			return err
		}
	}

	DB = db // передача открытого подключения к БД из временной локальной переменной в постоянную глобальную
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
