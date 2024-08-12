# Use the official Golang image as the base image
FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Set Go proxy to use proxy.golang.org and fall back to direct
ENV GOPROXY=https://proxy.golang.org,direct

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
