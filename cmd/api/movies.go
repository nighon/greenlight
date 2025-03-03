package main

import (
	"net/http"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.Movies.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
