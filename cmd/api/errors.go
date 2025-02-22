package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "there was a problem")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: %s path: %s err: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) unauthorized(w http.ResponseWriter, r *http.Request) {
	log.Printf("unauthorized: %s path: %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusUnauthorized, "you are not authorized to access this resource")
}

func (app *application) forbidden(w http.ResponseWriter, r *http.Request) {
	log.Printf("forbidden: %s path: %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "you do not have permission to access this resource")
}
