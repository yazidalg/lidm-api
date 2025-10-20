# --- Stage 1: Build the Go binary ---
# Use a specific Go version for consistency.
FROM golang:1.24.4-alpine3.21 AS builder

# Set the working directory inside the container.
WORKDIR /app

# Install git only if you have private modules, otherwise it's not needed.
# ca-certificates are needed for 'go mod download' to talk to HTTPS servers.
RUN apk add --no-cache ca-certificates

# Copy dependency files and download them. This is cached.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code.
COPY . .

# Build the application, creating a static binary.
# The -trimpath flag removes local filesystem paths from the binary for reproducibility.
# The -ldflags strip debug info to reduce size.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/main ./cmd


# --- Stage 2: Create the final, minimal, and secure image ---
# Use a specific, minimal base image.
FROM alpine:3.22.0

# Install the Certificate Authority certificates. This is the only package needed.
RUN apk add --no-cache ca-certificates tzdata

# Create a dedicated, non-privileged user and group for the application.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Switch to the non-privileged user.
USER appuser

# Set the working directory.
WORKDIR /app

# Copy ONLY the compiled binary and any necessary static assets from the builder stage.
COPY --from=builder /app/main .

# Cloud Run will set the PORT environment variable at runtime.
EXPOSE 8080

# The command to run your application.
CMD ["./main"]
