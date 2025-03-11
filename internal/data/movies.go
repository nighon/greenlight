package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`         // Unique identifier for the movie
	CreatedAt time.Time `json:"created_at"` // Time when the movie was added to our db
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title, year)
		VALUES (?, ?)`

	args := []any{movie.Title, movie.Year}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query = `SELECT id, created_at, title, year FROM movies WHERE id = ?`
	if err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year); err != nil {
		return err
	}

	return nil
}

func (m MovieModel) GetAll(ctx context.Context, title string) ([]*Movie, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, year
		FROM movies
		WHERE title like "%%%s%%"`, title)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// totalRecords := 0
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
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
