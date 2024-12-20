it:
	go mod tidy

cli:
	go run cmd/sentry-tunnel/main.go --help --trusted-sentry-dsn "http://host.docker.internal:8081/0"

.PHONY: sample
sample:
	go run sample/sample.go

.PHONY: benchmark
benchmark:
	wrk -t12 -c400 -d30s -s benchmarks/envelope.lua http://localhost:8080/tunnel
