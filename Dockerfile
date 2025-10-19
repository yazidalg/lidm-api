FROM golang:1.24.4-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-s -w" -o main ./cmd

# ✅ Ganti dari `scratch` ke `alpine` dan bawa certs
FROM alpine:3.22.0

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Opsional: minimal tambahan CA tools
RUN apk add --no-cache ca-certificates

EXPORT 3000

CMD ["./main"]
