FROM golang:1.12.0-alpine3.9 AS build

# create new user 
RUN adduser -D -g '' -u 2912 appuser

# Install git, dep and ca-certs
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN go get github.com/golang/dep/cmd/dep

# Set working directory
WORKDIR /go/src/go-echo

# Copy the source code
COPY main.go .

# Optimize build by removing debug informations and compile only for linux target and disabling cross compilation.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s -extldflags=-Wl,-z,now,-z,relro" -o /go/bin/go-echo .

# Take a new empty image
FROM scratch

# Copy the binary from the build image
COPY --from=build /go/bin/go-echo /go-echo

# Copy the users file from the build image
COPY --from=build /etc/passwd /etc/passwd

# Copy the ca certifications from the build image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Switch user
USER 2912

# Set working directory
WORKDIR /

# Run app
ENTRYPOINT ["/go-echo"]