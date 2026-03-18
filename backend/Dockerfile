# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
# go.sum is generated on first run; GONOSUMDB=* skips the checksum DB
COPY go.mod ./
ENV GONOSUMDB=*
ENV GOFLAGS=-mod=mod
RUN touch go.sum && go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
