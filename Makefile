it:
	docker buildx bake --load dev

run:
	docker run -it --rm -p 3003:3003 socheatsok78/sentry-tunnel:dev
