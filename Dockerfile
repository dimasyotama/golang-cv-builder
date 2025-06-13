# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cv-builder

# Final stage
FROM alpine:latest

# Install Chrome dependencies
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    ttf-freefont \
    fontconfig \
    dbus \
    xvfb

# Set Chrome binary location
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cv-builder .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./cv-builder"] 