# Use the official Go image to build the application
FROM golang:1.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the entire project to the working directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o concurrent-port-scanner .

# Use a smaller base image to run the app
FROM gcr.io/distroless/base-debian11

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/concurrent-port-scanner .

# Copy the internal/templates directory from builder
COPY --from=builder /app/internal/templates /app/internal/templates

# Expose port 8080 for the application
EXPOSE 8080

# Run the application
CMD ["./concurrent-port-scanner"]