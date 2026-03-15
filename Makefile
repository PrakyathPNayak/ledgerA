.PHONY: dev build test lint fmt migrate-up migrate-down migrate-create docs deploy-build

## dev: start docker-compose services (postgres + pgadmin) in background
dev:
	docker-compose up -d

## build: compile the Go server binary to bin/server
build:
	go build -o bin/server ./cmd/server

## test: run all Go tests with race detector
test:
	go test ./... -race -count=1

## lint: run golangci-lint on all source packages
lint:
	golangci-lint run ./cmd/... ./internal/... ./pkg/...

## fmt: format all Go source files with gofmt
fmt:
	gofmt -s -w ./cmd ./internal ./pkg

## migrate-up: apply all pending database migrations
migrate-up:
	migrate -path ./migrations -database "$${DATABASE_URL}" up

## migrate-down: roll back the most recent migration
migrate-down:
	migrate -path ./migrations -database "$${DATABASE_URL}" down

## migrate-create: scaffold a new named migration pair (usage: make migrate-create NAME=add_foo)
migrate-create:
	migrate create -ext sql -dir migrations -seq $${NAME}

## docs: list all generated documentation files
docs:
	@echo "Documentation files:"
	@find docs -type f -name '*.md' | sort
	@echo ""
	@echo "Root docs:"
	@ls -1 *.md 2>/dev/null || true

## deploy-build: build Linux amd64 binary for production deployment
deploy-build:
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 ./cmd/server
	@echo "Built: bin/server-linux-amd64"
