.PHONY: test run install

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go

run-docker:
	docker stop electrotech || true
	docker rm electrotech-back || true
	docker build -t electrotech-back .
	docker run -d -p 8080:8080 --name electrotech electrotech-back
