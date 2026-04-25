# Use stable Go version
FROM docker.io/library/golang:1.23



# Set working directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code and migrations
COPY . .
COPY migrations ./migrations

# Build the server binary at build time
RUN go build -o server ./

# Expose backend port
EXPOSE 8081

# Run server
CMD ["./server"]
