package main

import (
	"net/http"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthCheckResponse{
		Status:  "ok",
		Env:     app.config.env,
		Version: version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
