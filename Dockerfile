# image for building the server
FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# download dependency using go mod
RUN go mod download

# Build the application
RUN go build -o server ./

# Build a small image
FROM alpine:latest

# copy binary to new image
COPY --from=builder /build/server /

# Command to run
ENTRYPOINT ["/server"]