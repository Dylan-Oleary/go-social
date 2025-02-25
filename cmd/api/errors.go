package main

import (
	"log"
	"net/http"
)

func (app *application) badRequestError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) notFoundError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusNotFound, err.Error())
}
