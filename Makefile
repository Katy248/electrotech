.PHONY: test run install

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go

install:
	cp ./electrotech-back.service /lib/systemd/system/
