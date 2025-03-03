package main

import (
	"fmt"
	"net/http"

	"example.com/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		Year  int32  `json:"year"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title: input.Title,
		Year:  input.Year,
	}

	if err := app.models.Movies.Insert(movie); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

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
