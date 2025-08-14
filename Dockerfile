# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubevirt-mcp-server .

# Final stage
FROM registry.fedoraproject.org/fedora-minimal:latest

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/kubevirt-mcp-server .

# Expose the default MCP server port
EXPOSE 8080

# Run the binary
CMD ["./kubevirt-mcp-server"]