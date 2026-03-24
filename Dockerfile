# Build image with dependencies
FROM golang:1.25-alpine AS builder

# Current working directory
WORKDIR /app

# Copy dependency files and download before remainder of source code (layer caching)
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api/main.go

# Final stage with cleaner image
FROM alpine:3.19

WORKDIR /app

# Copy only the compiled binary from the build stage
COPY --from=builder /app/main .

# DEVTEST - MH
CMD ["echo", "Hello World!"]

# Listen for HTTP traffic
#EXPOSE 8080

# Runs compiled binary from line 20 - MH
#CMD ["./main"]