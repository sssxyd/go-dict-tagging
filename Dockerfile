# Use the official Golang image as the base image
FROM golang:1.22.6

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o app .

# Expose port 8080
EXPOSE 8080

# Set the entry point for the container
CMD ["./app"]