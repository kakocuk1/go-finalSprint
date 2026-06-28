package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kakocuk1/go-finalSprint/pkg/api"
)

const (
	webDir      = "./web"
	defaultPort = "7540"
)

func Run() error {
	port := os.Getenv("TODO_PORT") // реализована возмжность определить порт через переменную окружения(задание со звездочкой)
	if port == "" {
		port = defaultPort
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
