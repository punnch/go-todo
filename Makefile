include .env
export

export PROJECT_ROOT=$(shell pwd)

# env
env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "Do you want to clean up your environment files? You may lose your data. [y/N]: " ans; \
	if [ "$$ans" == "y" ]; then \
		docker compose down postgres && \
		sudo rm -rf out/pgdata && \
		echo "Environment files cleaned up successfuly"; \
	else \
		echo "Cleaning of environment files canceled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

# services
services-up: 
	@docker compose up -d db-service api-service

services-down:
	@docker compose down db-service api-service

services-rebuild:
	@docker compose build --no-cache db-service api-service

# migrate
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq not set. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; 
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action not set. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

# other
ps:
	@docker compose ps
