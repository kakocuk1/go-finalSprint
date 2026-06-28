package api

import (
	"net/http"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search != "" {
		tasks, err = db.TasksSearch(50, search) // 50 ограничение по заданию.
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeError(w, "Tasks get error: "+err.Error())
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
