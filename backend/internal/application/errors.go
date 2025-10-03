package application

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, code string, message string, details any) {
	data := map[string]any{
		"code":    code,
		"message": message,
		"details": details,
	}

	if err := writeJSON(w, status, data, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}

func (app *Application) routeNotFound(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "ROUTE_NOT_FOUND", "the requested route could not be found", nil)
}

func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "the requested method is not allowed for the specified route", nil)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors error) {
	validationErrors := errors.(validator.ValidationErrors)
	mappedErrors := make(map[string]string)
	for _, fieldError := range validationErrors {
		mappedErrors[fieldError.Field()] = fmt.Sprintf("failed on the '%s' tag", fieldError.Tag())
	}

	app.errorResponse(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "one or more fields failed validation", mappedErrors)
}
