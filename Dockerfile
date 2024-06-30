# Use the official Golang image as a base image
FROM golang:1.22-alpine as base

# Set the working directory in the container
WORKDIR /app

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Make port 5000 available to the world outside this container
EXPOSE 5000

# Run the Go application
CMD ["./main"]
