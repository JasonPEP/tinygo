# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o tinygo ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/tinygo .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./tinygo"]
