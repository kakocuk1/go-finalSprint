package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

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
		writeError(w, "JSON decode error: "+err.Error())
		return
	}

	if task.Title == "" {
		writeError(w, "Title is required")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
