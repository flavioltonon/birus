build:
	go build -o ./bin/birus ./cmd/api

check:
	go vet ./...

start: stop
	docker-compose up

stop:
	docker-compose down --remove-orphans

tidy:
	go fmt ./...