# Stage 1: Build the Go app
FROM golang:1.22 AS builder
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Change to the directory where the main.go file is located
WORKDIR /app/cmd/movie-api
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Build a small image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/movie-api/main .
# Expose port 9000 to the outside world
EXPOSE 9000
# Command to run the executable
CMD ["./main"]