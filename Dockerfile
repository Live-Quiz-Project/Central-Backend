# Use a smaller Alpine base image
FROM golang:alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

# Copy only necessary files for dependency download
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application
COPY . .

# Build the Go application
RUN go build -o app ./cmd/main.go

# Use a minimal base image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the port
EXPOSE 8080

# Set environment variable if needed
ENV USE_ENV_FILE=FALSE

# Run the application
CMD ["./app"]
