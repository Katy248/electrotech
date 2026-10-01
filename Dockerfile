FROM golang:latest AS builder

ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build


WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/srv ./cmd/server/main.go

FROM scratch

WORKDIR /bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/srv ./

ENV MIGRATIONS_DIR=/usr/share/srv/migrations
COPY --from=builder /app/sql/migrations /usr/share/srv/migrations

COPY --from=builder /app/electrotech-back.toml /etc/electrotech-back.toml

EXPOSE 8080 8021 30000-30020
ENTRYPOINT ["/bin/srv"]
