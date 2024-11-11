Build Stage
FROM golang:alpine3.19 AS BuildStage

WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . .

# Download and install any required dependencies
RUN go mod download

EXPOSE 80

# Build the Go app
RUN go build -o server -tags netgo main.go 

# Deploy Stage
FROM alpine:latest

WORKDIR /app

COPY --from=BuildStage /app/server .
# COPY --from=BuildStage /app/.env .

EXPOSE 80

# Command to run the executable
CMD ["./server"]
