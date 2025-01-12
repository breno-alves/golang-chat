FROM golang:1.23 AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -v -o server /app/cmd/chat/main.go


FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

CMD ["/app/server"]