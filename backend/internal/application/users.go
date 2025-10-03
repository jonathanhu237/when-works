package application

import (
	"errors"
	"net/http"

	"github.com/jonathanhu237/when-works/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(input); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	// Generate a random password
	password := app.generatePassword()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create user
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: string(passwordHash),
		IsAdmin:      false,
	}

	if err := app.models.User.Insert(user); err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameConflict):
			app.errorResponse(w, r, http.StatusConflict, "USER_USERNAME_CONFLICT", "username already exists", nil)
		case errors.Is(err, models.ErrEmailConflict):
			app.errorResponse(w, r, http.StatusConflict, "USER_EMAIL_CONFLICT", "email already exists", nil)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send email in background
	app.background(func() {
		data := map[string]any{
			"name":     user.Name,
			"username": user.Username,
			"password": password,
		}

		if err := app.mailer.SendHTML(user.Email, "Welcome to WhenWorks", "welcome.html", data); err != nil {
			app.logger.Error("failed to send welcome email to new user", "error", err)
			return
		}
		app.logger.Info("welcome email sent to new user", "email", user.Email)
	})

	// Return created user
	if err := app.writeJSON(w, http.StatusCreated, user, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
