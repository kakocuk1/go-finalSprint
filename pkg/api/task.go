package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

// Перенес addTask в общий task.go
// checkDate проверяет дату задачи и при необходимости корректирует её. Если дата пустая, ставит текущую дату.
// Если дата в прошлом, ставит текущую дату. Если дата в будущем, оставляет как есть. Если есть повторение, вычисляет следующую дату.
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date) // если неверный формат даты, возвращаем ошибку
	if err != nil {
		return err
	}
	// nowDate -текущая дата без времени, чтобы сравнивать только даты без учета времени
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// если есть повторение, вычисляем следующую дату
	// странно что по заданию, если нет правила повторения то ставим сегодняшнюю дату
	// то есть все мои прошлые активности пишутся сегодняшним числом ¯\_(ツ)_/¯
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if afterNow(nowDate, t) {
			task.Date = next
		}
	} else if afterNow(nowDate, t) {
		task.Date = now.Format(dateFormat)
	}

	return nil
}

// addTaskHandler - для HTTP
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "JSON decode error: "+err.Error())
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID is not specified")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	writeJson(w, task)
}

// editTaskHandler практически идентичен addTaskHandler, за исключением того,
// что проверяет еще task.ID и вызывает db.UpdateTask вместо db.AddTask.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "JSON decode error: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "ID is not specified")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is not specified")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJson(w, map[string]string{}) // Возвращаем пустой JSON-объект
}
