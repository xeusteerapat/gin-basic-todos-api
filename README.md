# Simple REST API

build with id flags

```shell
go build \
        -idflags "-X main.buildCommit=`git rev-parse --short HEAD` \
        -X main.buildTime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
        -o app
        "
```

## Example `.env`

```shell
PORT=8080
SIGN=YOUR_SECRET_PASSWORD
DB_CONN=USER:YOUR_SECRET_PASSWORD@tcp(127.0.0.1:3306)/YOUR_DB_NAME?charset=utf8mb4&parseTime=True&loc=Local
```
