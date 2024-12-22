it:
	go mod tidy

test:
	go test -v ./...

run:
	go run cmd/sentry-tunnel/main.go --log-level=debug

build:
	go build -o bin/sentry-tunnel cmd/sentry-tunnel/main.go

.PHONY: sample
sample:
	go run sample/sample.go

.PHONY: benchmark
benchmark:
	wrk -t12 -c400 -d30s -s benchmarks/envelope.lua http://localhost:8080/tunnel
