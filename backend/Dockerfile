# Dockerfile for the Go backend application

# --- Builder Stage ---
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# CGO_ENABLED=0 is for preventing dynamically linked libraries
# -o /server is the output file
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# --- Runner Stage ---
FROM alpine:latest

# Add ca-certificates to make SSL calls
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"] 