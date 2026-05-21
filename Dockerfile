FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git curl

# Install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./build/practiq-be ./cmd/api/

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/build/practiq-be /app/practiq-be
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/entrypoint.sh /app/entrypoint.sh

WORKDIR /app

RUN chmod +x entrypoint.sh

EXPOSE 8083

ENTRYPOINT ["./entrypoint.sh"]
