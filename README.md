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
