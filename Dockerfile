# Use the official Golang image as the base image
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .
COPY .env .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o cursor-crash-backend .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install PostgreSQL client (optional, if you need psql or other tools)
RUN apk --no-cache add postgresql-client

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/cursor-crash-backend .
COPY --from=builder /app/.env .

# Expose the port your application will run on
EXPOSE 8080

# Command to run the application
CMD ["./cursor-crash-backend"]