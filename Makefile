.PHONY: build test test-integration lint run migrate migrate-up dev-up dev-down dev-logs
build:
	go build ./cmd/controlplane
test:
	go test -race ./...
test-integration:
	TEST_DATABASE_URL="$${TEST_DATABASE_URL:-postgres://baobab:baobab@localhost:5432/baobab_control_plane?sslmode=disable}" \
	SHARED_CONTRACTS_DIR="$${SHARED_CONTRACTS_DIR:-../shared}" \
	go test -race ./...
lint:
	go vet ./...
run:
	go run ./cmd/controlplane
migrate:
	go run ./cmd/migrate
migrate-up:
	migrate -path internal/store/postgres/migrations -database "$(DATABASE_URL)" up
dev-up:
	docker compose up -d --wait
dev-down:
	docker compose down
dev-logs:
	docker compose logs -f
