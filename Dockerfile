# 🏗 Stage 1: Build the Go Application
FROM golang:1.24 AS builder

WORKDIR /app

# Copy module files first (caching dependencies)
COPY ../go.mod ../go.sum ./

RUN go mod tidy

# Copy entire project
COPY ../ ./

# Build the Go application
RUN go build -o gpsd-user-mgmt ./src

# 🏗 Stage 2: Use Debian Bookworm (latest GLIBC version)
FROM debian:bookworm

WORKDIR /app

# Install necessary libraries including latest GLIBC
RUN apt-get update && apt-get install -y curl lsof net-tools ca-certificates libc6 iputils-ping

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gpsd-user-mgmt .

# Expose the application port
EXPOSE 5500

# Start the application
CMD ["./gpsd-user-mgmt"]
