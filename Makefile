include .env
export

# docker
deploy:
	docker compose up -d

undeploy:
	docker compose down

deploy-app:
	docker compose up -d application

deploy-postgres:
	docker compose up -d postgres

deploy-migrate:
	docker compose upd -d migrate

# migrate
migrate-up:
	migrate -path migrations -database ${MIGRATE_DB_URL} up

migrate-down:
	migrate -path migrations -database ${MIGRATE_DB_URL} down

# other
service-run:
	go run cmd/goTodo/main.go
