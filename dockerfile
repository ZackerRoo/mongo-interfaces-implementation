# Use the official Golang image as a builder
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Set Go proxy
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o api .

# Use a minimal image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api .

# Expose the port the app runs on
EXPOSE 8085

ENV MONGO_URI="mongodb://root:123456789@mongo:27017/admin"

# Command to run the executable
CMD ["./api"]
