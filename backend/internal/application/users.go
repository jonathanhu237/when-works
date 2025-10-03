package application

import (
	"errors"
	"net/http"

	"github.com/jonathanhu237/when-works/backend/internal/models"
	"github.com/jonathanhu237/when-works/backend/internal/tasks"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required"`
	}

	if err := readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(input); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	// Generate a random password
	password := generatePassword()
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

	// Enqueue email task
	task, err := tasks.NewEmailNewUserTask(user.Email, user.Username, password, app.config)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if _, err := app.asynqClient.Enqueue(task); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return created user
	if err := writeJSON(w, http.StatusCreated, user, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
