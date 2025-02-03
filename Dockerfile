# # Start from the official Golang image
# FROM golang:1.18-alpine

# WORKDIR /app
# COPY . .

# # Install SQLite
# RUN apk update && apk add --no-cache sqlite sqlite-dev gcc musl-dev

# # Verify GCC installation
# RUN gcc --version

# EXPOSE 8080

# Stage 1: Build the Go application
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk update && apk add --no-cache gcc musl-dev sqlite sqlite-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final image
FROM alpine:3.16

WORKDIR /app

# Install SQLite
RUN apk update && apk add --no-cache sqlite

# Copy the built application and database from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/scans.db .

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["./main"]