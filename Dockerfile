# Stage 1: Build the application
FROM golang:1.21.6 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Run tests
FROM builder AS tester

# Change to directory with test files
WORKDIR /app/modules

# Run tests for each module
RUN go test ./payment
RUN go test ./donation
RUN go test ./user
RUN go test ./campaign

# Stage 3: Final image
FROM golang:1.21.6

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
