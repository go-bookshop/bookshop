include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## run/api/image: run the cmd/api application from bookshop docker image
.PHONY: run/api/image
run/api/image:
	docker run -p 4000:4000 --rm -v ./.env:/.env ${DOCKER_HUB_USERNAME}/bookshop:latest

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format all go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

## install/linters: install required linters
.PHONY: install/linters
install/linters:
	@echo 'Installing golangci-lint'
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

## tidy/diff: run mod tidy -diff and mod verify
.PHONY: tidy/diff
tidy/diff:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify

## test: run test with -race
.PHONY: test
test:
	@echo 'Running tests...'
	go test -race -vet=off ./...

## lint: run golangci-lint
.PHONY: lint
lint:
	@echo 'Running linters with golangci-lint...'
	golangci-lint run

## audit: run quality control checks
.PHONY: audit
audit:
	make -s tidy/diff
	make -s lint
	make -s test

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/api/image: build the docker image of cmd/api application
.PHONY: build/api/image
build/api/image:
	@echo 'Building cmd/api image...'
	docker build -t ${DOCKER_HUB_USERNAME}/bookshop:latest .

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## migrate/up: apply all migrations
.PHONY: migrate/up
migrate/up:
	@echo 'Applying migrations...'
	migrate -path ./migrations -database '${DATABASE_URL}' up

## migrate/down: rollback the last migration
.PHONY: migrate/down
migrate/down:
	@echo 'Rolling back the last migration...'
	migrate -path ./migrations -database '${DATABASE_URL}' down

## migrate/reset: rollback all migrations
.PHONY: migrate/reset
migrate/reset:
	@echo 'Rolling back all migrations...'
	migrate -path ./migrations -database '${DATABASE_URL}' down --all