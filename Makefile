SHELL := /bin/bash

.PHONY: dev-up dev-down test-backend api

dev-up:
	docker compose -f ops/docker-compose.yml up -d

dev-down:
	docker compose -f ops/docker-compose.yml down -v

test-backend:
	cd backend && go test ./...

api:
	cd backend && go run ./cmd/api
