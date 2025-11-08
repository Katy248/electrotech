FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/srv ./cmd/server/main.go
EXPOSE 8080
EXPOSE 8021
ENTRYPOINT ["/app/srv"]
