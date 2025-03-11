# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/api

# Stage 2: Use a minimal and secure runtime image
# FROM gcr.io/distroless/base
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Run the application with a non-root user (security best practice)
USER nonroot:nonroot

# Start the application
CMD ["./main"]
