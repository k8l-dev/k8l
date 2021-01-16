FROM golang:alpine as builder

RUN apk add --no-cache gcc musl-dev
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -tags "fts5" -a -v -o k8l .

FROM alpine:3.11.3
COPY --from=builder /build/k8l .

# executable
# ENTRYPOINT [ "./k8l" ]
# arguments that can be overridden
CMD ./k8l -listen $LISTEN -dbpath $DB