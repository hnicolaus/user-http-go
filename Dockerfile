# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.19-alpine as Build

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

# Generate RSA keys, copy them to the root directory, and build our binary at root location.
RUN apk add --no-cache openssl && \
    openssl genrsa -out /tmp/rsa 4096 && \
    openssl rsa -in /tmp/rsa -pubout -out /tmp/rsa.pub && \
    mv /tmp/rsa /tmp/rsa.pub / && \
    GOPATH= go build -o /main cmd/main.go

# This is the actual image that we will be using in production.
FROM alpine:latest

# We need to copy the binary from the build image to the production image.
COPY --from=Build /main .
COPY --from=Build /rsa /rsa.pub /

# This is the port that our application will be listening on.
EXPOSE 1323

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./main"]