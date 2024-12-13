.PHONY: up down build fresh logs test lint swaggo

up:
	@if [ ! -f .env ]; then \
        cp .env.local .env; \
    fi
	docker compose -f docker-compose.local.yaml up -d

down:
	docker compose -f docker-compose.local.yaml down && docker image rm project-api-project-http-api

build:
	docker build -t project-http-api:latest .

fresh:
	@if [ ! -f .env ]; then \
        cp .env.test .env; \
    fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V
	docker compose exec project-api go run .

migrate:
	docker compose exec project-api go run .

logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...

lint:
	golangci-lint run -v

swaggo:
	swag init -g **/**/*.go
