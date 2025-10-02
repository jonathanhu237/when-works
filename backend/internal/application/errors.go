package application

import "net/http"

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
	}

	if details != nil {
		data["details"] = details
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
