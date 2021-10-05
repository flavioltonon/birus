build:
	env GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
	go build -ldflags '-s -w' -o ./bin/birus ./cmd/api

check:
	go vet ./...

start: stop
	docker-compose up --build

stop:
	docker-compose down --remove-orphans

tidy:
	go fmt ./...