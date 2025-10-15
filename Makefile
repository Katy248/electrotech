.PHONY: test run install

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go

run:
	docker build -t electrotech-back .
	docker run -d -p 8080:8080 electrotech-back
