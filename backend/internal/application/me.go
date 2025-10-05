package application

import (
	"errors"
	"net/http"

	"github.com/jonathanhu237/when-works/backend/internal/models"
)

func (app *Application) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	requester := r.Context().Value(requesterContextKey).(*RequesterInfo)

	user, err := app.models.User.GetByID(requester.UserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found", nil)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]any{"user": user}, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
