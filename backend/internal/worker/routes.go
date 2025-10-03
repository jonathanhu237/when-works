package worker

import (
	"github.com/hibiken/asynq"
	"github.com/jonathanhu237/when-works/backend/internal/tasks"
)

func (w *Worker) routes() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeEmailNewUser, w.HandleEmailNewUserTask)

	return mux
}
