# Start from the official Golang image
FROM golang:1.18-alpine

WORKDIR /app

# Copy the source code
COPY . .

# Install SQLite and other dependencies
RUN apk update && apk add --no-cache sqlite sqlite-dev gcc musl-dev

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["sh", "-c", "go run main.go"]