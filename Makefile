# Makefile for user-service

.PHONY: run migrate-create migrate-up migrate-up-to migrate-up-by-one migrate-down migrate-down-to migrate-status migrate-version

# Run the application
run:
	go run cmd/api/main.go

# Create a new migration file
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<table-name>"; exit 1; fi
	goose create $(name) sql

# Apply all available migrations
migrate-up:
	goose up

# Migrate up to a specific version
migrate-up-to:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-up-to version=<version>"; exit 1; fi
	goose up-to $(version)

# Migrate up by one step
migrate-up-by-one:
	goose up-by-one

# Roll back a single migration
migrate-down:
	goose down

# Roll back to a specific version
migrate-down-to:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-down-to version=<version>"; exit 1; fi
	goose down-to $(version)

# Show migration status
migrate-status:
	goose status

# Show current database version
migrate-version:
	goose version
