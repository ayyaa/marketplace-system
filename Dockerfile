# Start with the official Go image
FROM golang:1.22-alpine as Build

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download and cache the Go modules dependencies
RUN go mod download

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the application source code
COPY . .

RUN swag init
# Build the Go app
RUN go build -o /main main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/main"]
