package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()
	router.NotFound(app.routeNotFound)
	router.MethodNotAllowed(app.methodNotAllowed)

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	return router
}
