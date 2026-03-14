include .env
export

# docker
deploy:
	docker compose up -d

undeploy:
	docker compose down

# local
run-api-service:
	go run cmd/api-service/main.go
run-db-service:
	go run cmd/db-service/main.go

# migrate
migrate-up:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} up

migrate-down:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} down

migrate-up-step:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} up ${step}

migrate-down-step:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} down ${step}