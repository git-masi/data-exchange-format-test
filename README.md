# Data Exchange Format Performance Test

## Server

Start the server with `go run cmd/server/main.go`.

Optionally set the Prometheus buckets using the `--buckets` flag. For example `go run cmd/server/main.go --buckets=js`.

## Client

### Go

Start the Go client with `go run cmd/client/main.go`

Optionally configure the type of message to send with the `--type` flag. For example `go run cmd/client/main.go --type=json`.

Optionally configure the WebSocket host name with the `--host` flag. For example `go run cmd/client/main.go --host="example.com"`.

### JavaScript

Start the JavaScript client with `cd cmd/client/ && npm start`.

Optionally configure the type of message to send by passing in an argument. For example `cd cmd/client/ && npm start json`.

Optionally configure the WebSocket host name by passing in an argument (note you must also pass in a message type argument). For example `cd cmd/client/ && npm start json "example.com"`.

## Dataviz

Start the data visualization app with `cd cmd/dataviz/ && go run *[^_test].go`.

Optionally use test mode to load charts with test data `cd cmd/dataviz/ && go run *[^_test].go --test=true`.

Optionally configure the URL to use to fetch Prometheus metrics with the `--url` flag. For example `cd cmd/dataviz/ && go run *[^_test].go --url="http://example.com/metrics"`
