# Makefile for managing migrations with golang-migrate

# Database connection parameters from environment variables
DB_USER := $(DB_USER)
DB_PASSWORD := $(DB_PASSWORD)
DB_HOST := $(DB_HOST)
DB_PORT := $(DB_PORT)
DB_NAME := $(DB_NAME)
DB_SSLMODE := $(DB_SSLMODE)

# The directory where your migrations are stored
MIGRATIONS_DIR := db/migrations

# The migrate command
MIGRATE_CMD := migrate

# Default goal is to show help
.DEFAULT_GOAL := help

# Run migrations up
up:
	@echo "Running migrations on ${DB_NAME}..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

# Rollback the last migration
down:
	@echo "Rolling back the last migration on ${DB_NAME}..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

# Check the current migration version
version:
	@echo "Checking migration version for ${DB_NAME}..."
	$(MIGRATE_CMD) -path $(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" version

# Show help for the Makefile
help:
	@echo "Usage:"
	@echo "  make up           # Apply migrations"
	@echo "  make down         # Rollback the last migration"
	@echo "  make version      # Check the current migration version"
	@echo "  make help         # Show this message"
