# Stage 1: Build the Go application
FROM golang:alpine AS builder

# Set CGO_ENABLED to 1 because this package requires CGO
ENV CGO_ENABLED=1

# Install necessary C libraries for go-sqlite3 and Alpine
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with static linking
RUN go build -ldflags='-s -w -extldflags "-static"' -o gopaster

# Stage 2: Create the runtime image
FROM alpine:latest

# Install necessary packages for runtime
RUN apk add --no-cache libc6-compat openssl

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/gopaster .

# Copy the templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Copy the startup script
COPY --from=builder /app/docker/entrypoint.sh .

# Mount volume for SSL certificates
VOLUME /app/certs

# Mount volume for SQLite database
VOLUME /app/data

# Expose the port the application runs on
EXPOSE 443

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Command to run the entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]