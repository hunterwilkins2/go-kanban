MAIN_PACKAGE_PATH := ./cmd/web
BINARY_NAME := go-kanban

## run: runs this package with hot reloading when saved
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0

## tailwind/build: complies tailwind css
.PHONY: tailwind/build
tailwind/build:
	tailwindcss -i ./ui/static/css/input.css -o ./ui/static/css/dist/output.css --minify

## build: builds this package
.PHONY: build
build: tailwind/build
	go build -o=bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## generate-sql: uses sqlc to generate sql models and queries
.PHONY: generate-sql
generate-sql:
	sqlite3 db/${BINARY_NAME}.db ".read db/schema.sql"
	sqlc -f db/sqlc.yaml generate

# Utilites
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'