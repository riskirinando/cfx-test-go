# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git
COPY go.mod ./
COPY go.su[m] ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && chmod 755 main

# Final stage
FROM alpine:latest

# Install required packages
RUN apk --no-cache add ca-certificates wget

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy binary and set ownership
COPY --from=builder /app/main .
RUN chmod 755 main && chown appuser:appgroup main

# Create a static directory for assets (optional)
RUN mkdir -p static && chown -R appuser:appgroup static

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use full path, no shell
CMD ["./main"]
