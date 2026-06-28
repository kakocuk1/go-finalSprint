package main

import (
	"log"
	"os"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
	"github.com/kakocuk1/go-finalSprint/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE") //  путь к БД берётся из TODO_DBFILE, если переменная задана, иначе используется scheduler.db
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
