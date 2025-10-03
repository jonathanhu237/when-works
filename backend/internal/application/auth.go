package application

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.validator.Struct(input); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	// Get user by username
	user, err := app.models.User.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(time.Duration(app.config.JWT.Expiration) * time.Second)
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(app.config.JWT.Secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set JWT as HttpOnly cookie
	cookie := &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  expirationTime,
		HttpOnly: app.config.Environment == config.Production,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	// Return
	if err = app.writeJSON(w, http.StatusOK, map[string]any{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the accessToken cookie
	cookie := &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: app.config.Environment == config.Production,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	// Return success response
	if err := app.writeJSON(w, http.StatusNoContent, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
