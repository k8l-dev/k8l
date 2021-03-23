FROM golang:alpine as builder

RUN apk add --no-cache gcc musl-dev libuv-dev linux-headers     git autoconf automake libc-dev  pkgconf libtool make sqlite-dev 
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN git clone https://github.com/canonical/raft.git && \
    cd raft && \
    autoreconf -i && \
    ./configure && \
    make && \
    make install  

RUN git clone https://github.com/canonical/dqlite.git && \
    cd dqlite && \
    autoreconf -i && \
   ./configure && \
   make && \
   make install

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY go go

# Build the application
RUN go build -tags libsqlite3,dqlite -a -o k8l ./go
# RUN go build -tags libsqlite3 -a -o k8l ./go

FROM alpine:3.11.3

RUN apk add sqlite-dev libuv-dev
COPY --from=builder /build/k8l /usr/bin/k8l
COPY --from=builder /usr/local/lib/libraft* /usr/lib/
COPY --from=builder /usr/local/lib/libdqlite* /usr/lib/
COPY ./static /static
ENTRYPOINT [ "k8l" ]