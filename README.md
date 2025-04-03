# Let's Go Further

## 开启服务

```console
$ open -a Docker        # 开启 Docker 服务
$ docker compose up -d  # 在 Docker 中运行 MySQL 数据库

$ source .envrc-example # 导入环境变量，设置参数
$ go run ./cmd/api      # 启动 Web 服务

$ go run ./cmd/api -port=3000 -env=production   # 或者从命令行设置参数，启动 Web 服务。
                                                # 优先级：命令行 > 环境变量 > 默认值
```

## 数据库迁移

安装迁移工具

```console
$ brew install golang-migrate
```

运行迁移脚本

```console
$ migrate -path=./migrations -database="mysql://dev:dev@tcp(127.0.0.1:3306)/greenlight" up
```

## 测试

```console
curl -i http://localhost:4000/v1/healthcheck
curl -i http://localhost:4000/v1/movies
curl -i http://localhost:4000/v1/movies/1
curl -i -X POST -d '{"title": "Terminator", "year": "1986"}' -- http://localhost:4000/v1/movies
curl -i -X POST -d '{"title": "Terminator 2", "year": "1991"}' -- http://localhost:4000/v1/movies
```
