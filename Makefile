it:
	docker buildx bake --load dev

run:
	docker run -it --rm -p 3003:3003 socheatsok78/sentry-tunnel:dev

.PHONY: benchmark
benchmark:
	wrk -t12 -c400 -d30s -s benchmarks/envelope.lua http://localhost:3003
