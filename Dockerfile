# Using Golang binaries as base image
FROM alpine:3.19 

# Current working directory
WORKDIR /app

# Copy project files into container
COPY go.mod go.sum ./

# Command to echo hello world
CMD ["echo", "Hello World!"]
