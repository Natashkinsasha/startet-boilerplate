.PHONY: run build wire migrate migrate-init migrate-rollback migrate-status migrate-create vet lint docker-up docker-down proto install-tools deps dev standalone swagger test test-unit test-integration test-functional

local-run:
	APP_ENV=local go run ./cmd/api/...
dev:
	air

build:
	go build -o bin/api ./cmd/api/...

wire:
	wire gen ./...

migrate:
	go run ./cmd/migrate/... migrate
migrate-init:
	go run ./cmd/migrate/... init
migrate-rollback:
	go run ./cmd/migrate/... rollback
migrate-status:
	go run ./cmd/migrate/... status
migrate-create:
	go run ./cmd/migrate/... create $(name)

test-unit:
	go test ./... -tags=unit -v -count=1 2>&1 | grep -v '\[no test files\]'

test-integration:
	go test ./... -tags=integration -v -count=1 2>&1 | grep -v '\[no test files\]'

test-functional:
	go test ./tests/functional/... -tags=functional -v -count=1 2>&1 | grep -v '\[no test files\]'

test:
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@$(MAKE) test-functional

vet:
	go vet -tags=$(shell grep -rh '//go:build' . --include='*.go' | awk '{print $$2}' | sort -u | paste -sd,) ./...

lint:
	golangci-lint run ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

install-tools:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest

deps:
	go mod tidy

swagger:
	APP_ENV=standalone go run ./cmd/swagger/...

proto:
	protoc -I proto \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		$(shell find proto -name "*.proto")