include .env

DOCKER = docker compose exec server
number :=

migrate-up:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path ./migrations up $(number)

migrate-down:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path ./migrations down $(number)

generate-api-doc:
	$(DOCKER) swag init