include .env
export

# docker
deploy:
	docker compose up -d

undeploy:
	docker compose down

# migrate
migrate-up:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} up

migrate-down:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} down

migrate-up-step:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} up ${step}

migrate-down-step:
	migrate -path internal/services/db-service/migrations -database ${LOCAL_DB_URL} down ${step}