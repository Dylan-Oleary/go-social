package main

import (
	"net/http"
)

func (app *application) badRequestError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal server Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) notFoundError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusNotFound, err.Error())
}
