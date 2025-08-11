# Dockerfile for building a Go application with distroless image

# Go version should match the version in go.mod
ARG GO_VERSION=1.23

# Use the official Go image as the build environment
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY app/ ./
RUN go mod tidy && \
    CGO_ENABLED=0 go build -o bin ./cmd/main.go

# Use a distroless image for the final build
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/bin ./bin

# Labels
LABEL maintainer="Corespark Engineering <oss@corespark.io>"
LABEL AUTHOR="Corespark Engineering <oss@corespark.io>"

# Set the entrypoint to the built binary
ENTRYPOINT ["/app/bin"]