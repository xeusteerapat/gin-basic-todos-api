# Simple REST API

build with id flags

```shell
go build \
        -idflags "-X main.buildCommit=`git rev-parse --short HEAD` \
        -X main.buildTime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
        -o app
        "
```
