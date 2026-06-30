FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the scraper
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scraper cmd/scraper/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/scraper .
COPY --from=builder /app/.env.example .env.example

# Expose no ports (this is a background worker)

# Run the scraper
CMD ["./scraper"]
