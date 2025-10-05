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

func (app *Application) UpdateMeHandler(w http.ResponseWriter, r *http.Request) {
	requester := r.Context().Value(requesterContextKey).(*RequesterInfo)

	var input struct {
		Email *string `json:"email" validate:"omitempty,email"`
		Name  *string `json:"name" validate:"omitempty,min=1"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Email == nil && input.Name == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	if err := app.validator.Struct(input); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

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

	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Name != nil {
		user.Name = *input.Name
	}

	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found", nil)
		case errors.Is(err, models.ErrEmailConflict):
			app.errorResponse(w, r, http.StatusConflict, "USER_EMAIL_CONFLICT", "email already exists", nil)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]any{"user": user}, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
