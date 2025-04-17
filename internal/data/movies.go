package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`                // Unique identifier for the movie
	CreatedAt time.Time `json:"-"`                 // Time when the movie was added to our db
	Title     string    `json:"title"`             // The title of the movie
	Year      int32     `json:"year,omitempty"`    // The release year of the movie
	// Runtime   Runtime   `json:"runtime,omitempty"` // The runtime of the movie in minutes
	// Genres    []string  `json:"genres,omitempty"`  // The genres of the movie
	// Version   int32     `json:"version"`           // The version of the movie: starts at 1 and increments each time the movie is updated
}

// func ValidateMovie(v *validator.Validator, movie *Movie) {
// 	v.Check(movie.Title != "", "title", "must be provided")
// 	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

// 	v.Check(movie.Year != 0, "year", "must be provided")
// 	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
// 	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

// 	v.Check(movie.Runtime != 0, "runtime", "must be provided")
// 	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

// 	v.Check(movie.Genres != nil, "genres", "must be provided")
// 	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
// 	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
// 	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
// }

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(ctx context.Context, movie *Movie) error {
	query := `INSERT INTO movies (title, year)
		VALUES (?, ?)`

	args := []any{movie.Title, movie.Year}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	query = `SELECT id, created_at, title, year FROM movies WHERE id = ?`
	if err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year); err != nil {
		return err
	}

	return nil
}

func (m MovieModel) Get(ctx context.Context, id int64) (*Movie, error) {
	if id < 1 {
		return nil, errors.New("invalid id")
	}

	query := `SELECT id, created_at, title, year FROM movies WHERE id = ?`

	var movie Movie

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &movie, nil
}

// func (app *application) updateMovie

func (m MovieModel) GetAll(ctx context.Context, title string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, year
		FROM movies
		WHERE title like "%%%s%%"`, title)

	// args := []interface{}{title}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// rows, err := m.DB.QueryContext(ctx, query, args...)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			// &totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	Metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, Metadata, nil
}
