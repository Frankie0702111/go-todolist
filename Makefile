include .env

number :=
migrate-up:
			migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(${APP_HOST}:${DB_PORT})/${DB_NAME}" -path ./migrations up $(number)

migrate-down:
			migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(${APP_HOST}:${DB_PORT})/${DB_NAME}" -path ./migrations down $(number)

