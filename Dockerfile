FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Install git and build dependencies (needed for some Go packages)
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o main ./cmd

FROM scratch

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
