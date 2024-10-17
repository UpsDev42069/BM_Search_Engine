# Stage 1: Build the Go application
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules manifests and download dependencies first 
# to leverage Docker's caching mechanism
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create a lightweight image to deploy the Go application
FROM alpine:3.18

# Install ca-certificates and SQLite
RUN apk --no-cache add ca-certificates sqlite

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# If you need to include .env, uncomment the line below
# COPY .env .

# Expose the application's port (ensure it matches the one specified in main.go)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]