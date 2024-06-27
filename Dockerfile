# Use the official Golang image as the build environment
FROM golang:1.22.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o node-affinity-guard .

# Use a minimal base image for the final build
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/node-affinity-guard /node-affinity-guard

# Command to run the binary
ENTRYPOINT ["/node-affinity-guard"]
