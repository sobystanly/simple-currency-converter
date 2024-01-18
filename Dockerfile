# Use a minimal base image with only necessary dependencies
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy only the necessary files for dependency resolution
COPY /src/go.mod /src/go.sum ./

# Download and cache Go dependencies
RUN go mod download

# Copy the entire application source code
COPY /src .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


# Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder image
COPY --from=builder /app/app .

# Expose the port the application runs on
EXPOSE 8001

# Run the application
CMD ["./app"]
