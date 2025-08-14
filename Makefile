BINARY_NAME=go-ecom

build:
	@go build -o $(BINARY_NAME).exe cmd/main.go

test:
	@go test -v ./...

run: build
	@./$(BINARY_NAME).exe

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down