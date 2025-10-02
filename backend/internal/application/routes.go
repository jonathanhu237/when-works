package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()

	return router
}
