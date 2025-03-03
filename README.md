# Let's Go Further

## 开启服务

```console
$ go run ./cmd/api -port=3000 -env=production
```

## 数据库迁移

安装迁移工具

```console
$ brew install golang-migrate
```

运行迁移脚本

```console
$ migrate -path=./migrations -database="mysql://dev:dev@tcp(127.0.0.1:3316)/greenlight" up
```

curl -i -X POST -d @movie.json -- http://localhost:4000/v1/movies

Got error when inserting data: "Error 1054 (42S22): Unknown column '$1' in 'field list'"

```go
func (m MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title, year)
		VALUES ($1, $2)`

	args := []any{movie.Title, movie.Year}

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
```
---

Is `m.DB.Exec()` a blocking operation? If it is, I think the two lines: `ctx, cancel := context.WithTimeout(); defer cancel()` is not necessary.

For your convinience, I attached more code snippets.


```go
package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"` // Unique identifier for the movie
	CreatedAt time.Time `json:"-"`  // Time when the movie was added to our db
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

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

```

---

Well, here is the original code for Postgres. Since I use MySQL instead, I modified the code for MySQL. That's why you see the `context` in used in my own code.

```go
func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Need to use QueryRow(Context) because of the RETURNING clause (which returns the id, created_at and version)
	// RETURNING is a Postgres feature that is not part of SQL standard
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

```

Thanks! Could you modify the method "GetAll" for MySQL as well?

You can ignore the fields: "runtime, genres, version", as I don't have these fields currently.


```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres && $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}
```

---


```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	// Construct the WHERE clause
	query := fmt.Sprintf(`
		SELECT SQL_CALC_FOUND_ROWS id, created_at, title, year
		FROM movies
		WHERE (MATCH(title) AGAINST(?) OR ? = '')
		%s
		ORDER BY %s %s, id ASC
		LIMIT ? OFFSET ?`, genreFilterClause(len(genres)), filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare query arguments
	args := []any{title, title}
	args = append(args, genreArgs(genres)...) // Add genre conditions dynamically
	args = append(args, filters.limit(), filters.offset())

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Fetch total record count (alternative to count(*) OVER() in PostgreSQL)
	var totalRecords int
	err = m.DB.QueryRowContext(ctx, "SELECT FOUND_ROWS()").Scan(&totalRecords)
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return movies, metadata, nil
}

```

---

```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	// Base SELECT query
	query := fmt.Sprintf(`
		SELECT id, created_at, title, year
		FROM movies
		WHERE (MATCH(title) AGAINST(?) OR ? = '')
		%s
		ORDER BY %s %s, id ASC
		LIMIT ? OFFSET ?`, genreFilterClause(len(genres)), filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare query arguments
	args := []any{title, title}
	args = append(args, genreArgs(genres)...) // Add genre filter dynamically
	args = append(args, filters.limit(), filters.offset())

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year); err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// 🔹 Second query: Get total count of matching rows
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM movies
		WHERE (MATCH(title) AGAINST(?) OR ? = '')
		%s`, genreFilterClause(len(genres)))

	countArgs := []any{title, title}
	countArgs = append(countArgs, genreArgs(genres)...)

	var totalRecords int
	err = m.DB.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalRecords)
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return movies, metadata, nil
}


```
