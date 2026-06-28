package api

import (
	"net/http"
	"time"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "ID is not specified")
		return
	}

	task, err := db.GetTask(id) // получаем задачу по ID, чтобы проверить правило повтора
	if err != nil {
		writeError(w, "task not found")
		return
	}

	// если правило повтора не указано - удаляем задачу и вернем {}
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error())
			return
		}
		writeJson(w, map[string]string{})
		return
	}

	// если задача периодическая, то вычисляем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	// обновляем только дату
	if err := db.UpdateDate(id, next); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJson(w, map[string]string{})
}
