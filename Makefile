swagger:
	swag init -g ./cmd/main.go -o ./docs
local:
	POSTGRES_DB=example \
	POSTGRES_HOST=localhost \
	POSTGRES_USER=example \
	POSTGRES_PASSWORD=secret \
	POSTGRES_PORT=5433 \
	POSTGRES_SSLMODE=disable \
	HTTP_PORT=8080 \
	go run cmd/server/main.go

run:
