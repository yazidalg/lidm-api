FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o main ./cmd

# ✅ Ganti dari `scratch` ke `alpine` dan bawa certs
FROM alpine:3.21.3

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Opsional: minimal tambahan CA tools
RUN apk add --no-cache ca-certificates

CMD ["./main"]
