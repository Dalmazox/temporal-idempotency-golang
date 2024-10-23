DB_CONN ?= "postgres://docker:d0cker@localhost:5432/bank?sslmode=disable"

migrate:
	migrate -path ./migrations -database $(DB_CONN) up
.PHONY: migrate

migrate-down:
	migrate -path ./migrations -database $(DB_CONN) down
.PHONY: migrate-down