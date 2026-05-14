# Variables
DB_ADDR="postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"
MIGRATIONS_DIR=./cmd/migrate/migrations

# 1. Run all "UP" migrations
migrate-up:
	migrate -path=$(MIGRATIONS_DIR) -database=$(DB_ADDR) up

# 2. Run all "DOWN" migrations (Rollback)
migrate-down:
	migrate -path=$(MIGRATIONS_DIR) -database=$(DB_ADDR) down

# 3. Create a new migration file dynamically
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name."; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(name)

# 4. Force a specific migration version (used to fix dirty states)
migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: Please provide a version number."; \
		echo "Usage: make migrate-force version=4"; \
		exit 1; \
	fi
	migrate -path=$(MIGRATIONS_DIR) -database=$(DB_ADDR) force $(version)

seed:
	@go run cmd/migrate/seed/main.go

gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

test:
	@go test -v ./...