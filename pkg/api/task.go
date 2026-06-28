package api

import (
	"encoding/json"
	"net/http"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "ID is not specified")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Task not found")
		return
	}

	writeJson(w, task)
}

// editTaskHandler практически идентичен addTaskHandler, за исключением того,
// что проверяет еще task.ID и вызывает db.UpdateTask вместо db.AddTask.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "JSON decode error: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, "ID is not specified")
		return
	}

	if task.Title == "" {
		writeError(w, "Title is not specified")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJson(w, map[string]string{}) // Возвращаем пустой JSON-объект
}
