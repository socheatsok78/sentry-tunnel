FROM denoland/deno:2.1.4
# The port that your application listens to.
EXPOSE 3003

WORKDIR /app

# Prefer not to run as root.
USER deno

# Cache the dependencies as a layer (the following two steps are re-run only when deno.json & deno.lock is modified).
COPY deno.json deno.lock /app/
RUN deno install --frozen

# These steps will be re-run upon each file change in your working directory:
COPY main.ts /app
# Compile the main app so that it doesn't need to be compiled each startup/entry.
RUN deno cache main.ts

CMD ["run", "--allow-env", "--allow-net", "main.ts"]
