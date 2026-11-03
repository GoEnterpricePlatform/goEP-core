# Stage 1: Build the Go application
FROM golang:1.26 AS builder

WORKDIR /app

# Copy Go dependency files first to take advantage of Docker cache
COPY go.mod go.sum ./

# Download the dependencies declared in go.mod
RUN go mod download

# Copy the rest of the project files
COPY . .

# Move to the production application entrypoint
WORKDIR /app/cmd/api/env/prod

# Build a Linux production binary without CGO dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -tags=prod -o /app/goep-core


# Stage 2: Runtime image
FROM alpine:3.23

WORKDIR /app

# Install CA certificates for outbound HTTPS connections
RUN apk add --no-cache ca-certificates

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/goep-core /app/goep-core

# Copy runtime resources from pkg (.html, .md, etc.)
COPY --from=builder /app/pkg /app/pkg

# Copy the already-generated static files
COPY --from=builder /app/addons/web-app/resources/static /app/addons/web-app/resources/static

# Start the application
CMD ["./goep-core"]