package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) routes() http.Handler {
	router := chi.NewRouter()
	router.NotFound(app.notFound)
	router.MethodNotAllowed(app.methodNotAllowed)

	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Route("/v1/auth", func(r chi.Router) {
		r.Post("/login", app.LoginHandler)
		r.Post("/logout", app.LogoutHandler)
	})
	router.With(app.requireAuth).Route("/v1/me", func(r chi.Router) {
		r.Get("/", app.GetMeHandler)
		r.Patch("/", app.UpdateMeHandler)
		r.Post("/update-password", app.UpdateMePasswordHandler)
	})
	router.With(app.requireAuth, app.requireAdmin).Route("/v1/users", func(r chi.Router) {
		r.Get("/", app.ListUsersHandler)
		r.Post("/", app.CreateUserHandler)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", app.GetUserHandler)
			r.Patch("/", app.UpdateUserHandler)
			r.Delete("/", app.DeleteUserHandler)
			r.Post("/reset-password", app.ResetUserPasswordHandler)
		})
	})

	return router
}
