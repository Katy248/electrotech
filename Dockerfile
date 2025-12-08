FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/srv ./cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/srv /app/srv
EXPOSE 8080 8021 30000-30020
ENTRYPOINT ["/app/srv"]
