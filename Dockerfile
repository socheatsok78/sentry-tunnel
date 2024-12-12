FROM denoland/deno:alpine
EXPOSE 3003
WORKDIR /app
USER deno
COPY main.ts /app
CMD ["run", "--allow-env", "--allow-net", "main.ts"]
