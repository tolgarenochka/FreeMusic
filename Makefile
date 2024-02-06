run: update_swag build_app run_app

update_swag:
	swag init -g cmd/main.go

build_app: update_swag
	go build -o ./build/free_music_app ./cmd/main.go

run_app:
	./build/free_music_app

lint:
	golangci-lint run

test:
	go test -v ./...

COMPOSE_FILE = docker-compose.yml
up-mongobd:
	docker-compose -f $(COMPOSE_FILE) up -d

.PHONY: all build run