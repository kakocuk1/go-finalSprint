package api

import (
	"net/http"

	"github.com/kakocuk1/go-finalSprint/pkg/db"
)

// Правки по ревью - выношу лимит в константу
const tasksLimit = 50

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
		tasks, err = db.TasksSearch(tasksLimit, search)
	} else {
		tasks, err = db.Tasks(tasksLimit)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Tasks get error: "+err.Error())
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
