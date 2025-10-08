.PHONY: test run install

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go

install:
	go build -o ./electrotech-back cmd/server/main.go
	cp ./electrotech-back /usr/bin/
	cp ./electrotech-back.service /lib/systemd/system/
