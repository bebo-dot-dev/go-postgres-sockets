# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.16.7 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

RUN go get -d -v ./...

# Builds the application inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main

# Using a scratch stage - https://hub.docker.com/_/scratch/
FROM scratch AS runtime

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /main

# Run the web service on container startup.
CMD ["/main"]