.PHONY: build maria

build:
	go build \
		-ldflags "-X main.buildCommit=`git rev-parse --short HEAD` \
		-X main.buildTime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
		-o app

maria:
	docker run -p 127.0.0.1:3306:3306  --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=${MARIADB_ROOT_PASSWORD} -e MARIADB_DATABASE=myapp -d mariadb:latest

image:
	docker build -t gin-basic-todo-api:test -f Dockerfile .

container:
	docker run -p:8080:8080 --env-file ./local.env --link some-mariadb:db \
	--name myapp gin-basic-todo-api:test

installvegeta:
	go install github.com/tsenart/vegeta@latest
vegeta:
	echo "GET http://:8081/limit" | vegeta attack -rate=10/s -duration=1s | vegeta report
load:
	echo "GET http://:8081/limit" | vegeta attack -rate=10/s -duration=1s > results.10qps.bin
plot:
	 cat results.10qps.bin | vegeta plot > plot.10qps.html
hist:
	cat results.10qps.bin | vegeta report -type="hist[0,100ms,200ms,300ms]"