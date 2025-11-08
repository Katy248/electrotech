FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/srv ./cmd/server/main.go
EXPOSE 8080 8021 30000-30020
ENTRYPOINT ["/app/srv"]
