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
	router.Route("/v1/auth", func(r chi.Router) {
		r.Post("/login", app.LoginHandler)
	})

	return router
}
