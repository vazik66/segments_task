.PHONY: build run docs

include ./.env

build:
	go build -o ./build/ ./cmd/segment/main.go

run: build
	./build/main.exe

migrate:
	migrate -source file://pkg/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable up

downgrade:
	migrate -source file://pkg/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable down

docs:
	./scripts/docs
