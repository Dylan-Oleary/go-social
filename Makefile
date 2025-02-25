include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: dev
dev:
	@air

.PHONY: db-migrate-create
db-migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: db-migrate-up
db-migrate-up:
	@migrate -path=${MIGRATIONS_PATH} -database=${DB_ADDR} up

.PHONY: db-migrate-down
db-migrate-down:
	@migrate -path=${MIGRATIONS_PATH} -database=${DB_ADDR} down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: db-migrate-down-all
db-migrate-down-all:
	@migrate -path=${MIGRATIONS_PATH} -database=${DB_ADDR} down -all

.PHONY: db-reset
db-reset:
	@make db-migrate-down-all && make db-migrate-up && make db-seed

.PHONE: db-seed
db-seed:
	@go run ./cmd/migrate/seed/main.go