# Let's Go Further

## 开启服务

```console
# 开启 Docker 服务
$ open -a Docker

# 在 Docker 中运行 MySQL 数据库
$ docker compose up -d

# 启动 Web 服务
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

curl -i -X POST -d '{"title": "Terminator 2", "year": "1986"}' -- http://localhost:4000/v1/movies




















---

I want to return the inserted movie.

Please review my code:

```go
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

	query = `SELECT id, created_at, title, year FROM movies WHERE id = ?`
	args = []any{id}
	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year); err != nil {
		return err
	}

	return nil
}
```



---

I see the cancel() function has been replaced by the seond one. Will there be a problem?

```go
// Use a timeout for INSERT
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

// Use a NEW context for SELECT (prevents timeout issues from INSERT)
ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

