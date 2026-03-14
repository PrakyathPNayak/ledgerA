.PHONY: dev build test lint fmt migrate-up migrate-down

dev:
	docker-compose up -d

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -race -count=1

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

migrate-up:
	migrate -path ./migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path ./migrations -database "$${DATABASE_URL}" down
