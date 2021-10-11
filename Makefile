build:
	env GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
	go build -ldflags '-s -w' -o ./bin/birus ./cmd/api

check:
	go vet ./...	

image:
	docker build -t flavioltonon/birus:latest .

push:
	docker push flavioltonon/birus:latest

start: stop
	docker-compose up

stop:
	docker-compose down --remove-orphans

tidy:
	go fmt ./...