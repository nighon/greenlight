# Greenlight

## 开启服务

先开启 Docker 服务

```zsh
$ open -a Docker
```

在 Docker 中运行 Postgres 数据库 `# docker compose up -d --pull=never`

```zsh
$ source .envrc
$ docker compose up -d
$ go run ./cmd/api

# 或者从命令行设置参数，启动 Web 服务。
# 优先级：命令行 > 环境变量 > 默认值
$ go run ./cmd/api -port=3000 -env=production
```

## 数据库迁移

安装迁移工具

```zsh
$ brew install golang-migrate
```

运行迁移脚本

```zsh
$ migrate -path=./migrations -database="postgres://postgres:password@127.0.0.1:3306/greenlight" up
```

## 测试

```zsh
curl -i http://localhost:4000/v1/healthcheck
curl -i http://localhost:4000/v1/movies
curl -i http://localhost:4000/v1/movies/1
curl -i -X POST -d '{"title": "Terminator", "year": "1986"}' -- http://localhost:4000/v1/movies
curl -i -X POST -d '{"title": "Terminator 2", "year": "1991"}' -- http://localhost:4000/v1/movies
```
