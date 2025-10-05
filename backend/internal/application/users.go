package application

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		app.internalServerError(w, r, err)
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
			app.internalServerError(w, r, err)
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
	if err := app.writeJSON(w, http.StatusCreated, map[string]any{"user": user}, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.User.GetAll()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := map[string]any{"users": users}
	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		app.badRequestResponse(w, r, errors.New("user id is required"))
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid user id"))
		return
	}

	user, err := app.models.User.GetByID(userID)
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

func (app *Application) ResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		app.badRequestResponse(w, r, errors.New("user id is required"))
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid user id"))
		return
	}

	user, err := app.models.User.GetByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found", nil)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	password := app.generatePassword()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user.PasswordHash = string(passwordHash)

	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found", nil)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.background(func() {
		data := map[string]any{
			"name":     user.Name,
			"username": user.Username,
			"password": password,
		}

		if err := app.mailer.SendHTML(user.Email, "Your WhenWorks Password Was Reset", "password_reset.html", data); err != nil {
			app.logger.Error("failed to send password reset email", "error", err, "email", user.Email)
			return
		}
		app.logger.Info("password reset email sent", "email", user.Email)
	})

	if err := app.writeJSON(w, http.StatusNoContent, nil, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		app.badRequestResponse(w, r, errors.New("user id is required"))
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid user id"))
		return
	}

	var input struct {
		Email   *string `json:"email" validate:"omitempty,email"`
		Name    *string `json:"name" validate:"omitempty,min=1"`
		IsAdmin *bool   `json:"is_admin"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Email == nil && input.Name == nil && input.IsAdmin == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	if err := app.validator.Struct(input); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	user, err := app.models.User.GetByID(userID)
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
	if input.IsAdmin != nil {
		user.IsAdmin = *input.IsAdmin
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
