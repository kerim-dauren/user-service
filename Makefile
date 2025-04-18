.PHONY: godoc

godoc:
	swag init --dir ./cmd/app,./internal/api/http/v1,./internal/models --output ./internal/api/docs

.PHONY: compose-run
compose-run:
	docker compose down
	docker compose up --build --remove-orphans -d

.PHONY: go-test
go-test:
	go test ./... -v -cover

.PHONY: go-proto-gen
go-proto-gen:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	gen/proto/user.proto

.PHONY: generate-mocks
generate-mocks:
	mockgen -source=internal/domain/user.go -destination=internal/domain/mock_user.go -package=domain

.PHONY: migration-up migration-down migration-create

CURRENT_DIR := $(shell pwd)


MIGRATION_DIR := $(CURRENT_DIR)/db/migration/scripts

# make migration up - Run migration up
migration-up:
	goose -dir $(MIGRATION_DIR) postgres "$(PG_URL)" up

# make migration down - Run migration down
migration-down:
	goose -dir $(MIGRATION_DIR) postgres "$(PG_URL)" down

# make migration create - Create new migration sql file
migration-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATION_DIR) -s create "$$name" sql

.PHONY: project-build
# make project-build - Build project for linux amd64
project-build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o build/user-service cmd/main.go