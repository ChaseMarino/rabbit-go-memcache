# Use the official Golang image as the base image
FROM golang:1.17

# Set environment variables
ENV RABBITMQ_URL todo
ENV MEMCACHE_URL todo

# Set the working directory to /app
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o app .

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./app"]

