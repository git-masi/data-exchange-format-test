# Data Exchange Format Performance Test

## About

The purpose of this repo is to create a test that can answer the question: "does the data exchange format used between client and server impact performance?"

### This is for fun

Please note this is just for fun and is not meant to be academically rigorous. Nor is this meant to serve as a guide for what you should or should not due in a production application. The code found here may or may not be idiomatic, is probably not clean, and is likely to be very, very WET (meow). Many good coding practices have been thrown out the window in favor of getting something done quickly.

## Server

Start the server with `go run cmd/server/main.go`.

Optionally set the Prometheus buckets using the `--buckets` flag. For example `go run cmd/server/main.go --buckets=js`.

## Client

### Go

Start the Go client with `go run cmd/client/main.go`

Optionally configure the type of message to send with the `--type` flag. For example `go run cmd/client/main.go --type=json`.

Optionally configure the WebSocket host name with the `--host` flag. For example `go run cmd/client/main.go --host="example.com"`.

### JavaScript

Start the JavaScript client with `cd cmd/js-client/ && npm start`.

Optionally configure the type of message to send by passing in an argument. For example `cd cmd/js-client/ && npm start json`.

Optionally configure the WebSocket host name by passing in an argument (note you must also pass in a message type argument). For example `cd cmd/js-client/ && npm start json "example.com"`.

## Dataviz

Start the data visualization app with `cd cmd/dataviz/ && go run *[^_test].go`.

Optionally use test mode to load charts with test data `cd cmd/dataviz/ && go run *[^_test].go --test=true`.

Optionally configure the URL to use to fetch Prometheus metrics with the `--url` flag. For example `cd cmd/dataviz/ && go run *[^_test].go --url="http://example.com/metrics"`
