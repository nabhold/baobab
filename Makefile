.PHONY: build test lint run migrate-up
build:
	go build ./cmd/controlplane
test:
	go test -race ./...
lint:
	go vet ./...
run:
	go run ./cmd/controlplane
migrate-up:
	migrate -path internal/store/postgres/migrations -database "$(DATABASE_URL)" up
