# First stage: build the Go binary
FROM golang:1.23.8 AS builder

WORKDIR /app

# Copy Go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (consider giving the binary a more generic name)
RUN CGO_ENABLED=0 GOOS=linux go build -o xyz-credit-plus-be .

# Second stage: minimal runtime
FROM alpine:latest

# Install required packages
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/xyz-credit-plus-be /xyz-credit-plus-be

# Ensure the binary is executable
RUN chmod +x /xyz-credit-plus-be

# Make directory to save KYC documents
RUN mkdir -p /app/assets

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["/xyz-credit-plus-be"]