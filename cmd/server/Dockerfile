FROM golang:1.20.5

WORKDIR /usr/src/app

# Copy go.mod and go.sum from the build context
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the files from cmd/server
COPY cmd/server/ .
RUN go build -v -o /usr/local/bin/app ./main.go
