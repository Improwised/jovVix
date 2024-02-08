MIGRATE_BIN := $(shell which migrate)
SWAGGER_BIN := $(shell which swagger)
BIN := /usr/local/bin

.DEFAULT_GOAL := intro

intro:
	@echo "please specify a target {migrate, swagger-gen, start, start-api, migration-up, start-dev, test, clean-test-cache, test-wo-cache}"

migrate:
ifeq ($(MIGRATE_BIN),)
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz
	mv migrate.linux-amd64 migrate
	chmod +x migrate
	sudo mv ./migrate $(BIN)
endif
	migrate create -ext sql -dir database/migrations $(file_name)

swagger-gen:
ifeq ($(SWAGGER_BIN),)
	curl -L https://github.com/go-swagger/go-swagger/releases/download/v0.26.1/swagger_linux_amd64 --output swagger
	chmod +x swagger
	sudo mv ./swagger $(BIN)
endif
	swagger generate spec -o ./assets/swagger.json

start-api:
	go run app.go api

migrate-up:
	go run app.go migrate up

start-api-dev:
	@nodemon --exec go run app.go api --signal SIGTERM

test:
	go -v test ./...

clean-test-cache:
	go clean -testcache && go test ./...

# Test without cache
test-wo-cache: clean-test-cache test

build:
	go build -o=$(app_name) .

install:
	go build -ldflags="-s -w" -o=$(app_name) .
