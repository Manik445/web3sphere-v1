.PHONY: build run dev test lint migrate-up migrate-down migrate-create seed docker-up docker-down swagger clean

# Variables
APP_NAME=web3sphere
DB_URL=postgres://postgres:postgres@localhost:5432/web3sphere?sslmode=disable
MIGRATE_PATH=migrations

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

dev:
	air -c .air.toml

test:
	go test -v -cover ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $$name

seed:
	psql "$(DB_URL)" -f migrations/seeds/000001_seed_data.up.sql

seed-down:
	psql "$(DB_URL)" -f migrations/seeds/000001_seed_data.down.sql

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

swagger:
	swag init -g cmd/server/main.go -o docs/swagger

clean:
	rm -rf bin/
	go clean
