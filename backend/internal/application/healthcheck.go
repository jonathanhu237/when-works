package application

import "net/http"

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"environment": string(app.config.Environment),
		},
	}

	if err := writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
