.PHONY: test run

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go
