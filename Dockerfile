# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build-time dependencies
RUN apk add --no-cache git

# Copy and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /cv-builder .

# --- Final stage ---
FROM alpine:latest

# Install runtime dependencies for headless Chrome
RUN apk add --no-cache \
    chromium \
    ttf-freefont \
    dbus \
    fontconfig

# Create a non-root user for security
RUN adduser -D -g '' appuser

# Set the working directory for the final container.
# All subsequent commands (like CMD) will run from here.
WORKDIR /app

# Copy the compiled binary from the builder stage into the new WORKDIR
COPY --from=builder /cv-builder .

# Copy the templates directory into the new WORKDIR
COPY --from=builder /app/templates ./templates/

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 9090

# Since WORKDIR is /app, this command is now equivalent to /app/cv-builder
CMD ["./cv-builder"]