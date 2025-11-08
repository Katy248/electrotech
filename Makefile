.PHONY: test run install

test:
	go test ./...

run:
	sqlc generate
	go run cmd/server/main.go

IMAGE := electrotech-back
CONTAINER_NAME := electrotech

run-docker:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	docker build -t $(IMAGE) .
	docker run \
		--volume ./example:/data\
		--detach \
		-p 8080:8080 \
		-p 8021:8021 \
		--name $(CONTAINER_NAME) $(IMAGE)
