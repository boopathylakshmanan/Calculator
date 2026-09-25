DB_URL := postgres://calculator:calculator@postgres:5432/calculator?sslmode=disable
NETWORK := calculator_net

.PHONY: up down migrate test

up:
	docker compose up -d

down:
	docker compose down

migrate:
	docker run --rm -v "$(CURDIR)/migrations:/migrations" --network $(NETWORK) migrate/migrate -path=/migrations -database "$(DB_URL)" up

test:
	cd backend && go test ./...
	cd frontend && npm test
