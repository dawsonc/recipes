# Start from a Golang base image
FROM golang:alpine

# Set the current working directory inside the container
WORKDIR /recipes

# Copy the necessary files to the container
COPY .. .
RUN go mod download

# Build the application
RUN go build -o main ./app/*.go

# Expose port 8080 for the application
EXPOSE 8080

# Set the entry point for the container
CMD ["./main"]
