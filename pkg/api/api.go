package api

import (
	"encoding/json"
	"net/http"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

func Init() {
	InitAuth()
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}

// writeJson принимает любые данные, превращает их в JSON через json.Marshal, выставляет правильный заголовок Content-Type, пишет в ответ
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

// writeError- превращает текст ошибки в мапу с ключом "error", чтобы получился JSON: {"error": "..."}
func writeError(w http.ResponseWriter, status int, errText string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	resp, _ := json.Marshal(map[string]string{"error": errText})
	w.Write(resp)
}

// taskHandler- для HTTP
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method is not allowed")
	}
}

// deleteHandler используется для удаления задачи по ID.
// Если задача периодическая, то она удаляется полностью, без создания новой даты.
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID is not specified")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJson(w, map[string]string{})
}
