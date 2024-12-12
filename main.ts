//    _____            __                ______                       __
//   / ___/___  ____  / /________  __   /_  __/_  ______  ____  ___  / /
//   \__ \/ _ \/ __ \/ __/ ___/ / / /    / / / / / / __ \/ __ \/ _ \/ / 
//  ___/ /  __/ / / / /_/ /  / /_/ /    / / / /_/ / / / / / / /  __/ /  
// /____/\___/_/ /_/\__/_/   \__, /    /_/  \__,_/_/ /_/_/ /_/\___/_/   
//                          /____/                                      
// A tunnel is an HTTP endpoint that acts as a proxy between Sentry and your application.
// Learn more at https://docs.sentry.io/platforms/javascript/troubleshooting/#using-the-tunnel-option

import * as log from "@std/log";

type SentryHeaderStruct = {
  dsn: string;
  event_id: string;
  sent_at: string;
  sdk: {
    name: string;
    version: string;
  };
}


if (import.meta.main) {
  log.setup({
    handlers: {
      default: new log.ConsoleHandler("DEBUG", {
        formatter: log.formatters.jsonFormatter,
        useColors: false,
      }),
    },
  });

  const options: Deno.ServeTcpOptions = {
    port: parseInt(Deno.env.get("PORT") as string || "3003"),
    onListen(localAddr) {
      log.info(`Listening on http://localhost:${localAddr.port}`);
      log.info("Sentry tunnel started! Log data will stream in below:");
    },
  }

  Deno.serve(options, async (req) => {
    if (req.method !== "POST") {
      return new Response(null, { status: 405 });
    }

    try {
      const envelopeBytes = await req.arrayBuffer();
      const envelope = new TextDecoder().decode(envelopeBytes);
      const piece = envelope.split("\n")[0];
      const header: SentryHeaderStruct = JSON.parse(piece);

      if (!header.dsn) {
        log.error("No DSN found in the envelope header");
        return new Response(JSON.stringify({ error: "No DSN found in the envelope header" }), { status: 400 });
      }

      // Extract the project ID from the DSN
      const dsn = new URL(header.dsn);
      const project_id = dsn.pathname?.replace("/", "");

      if (!project_id) {
        log.error("No project ID found in the DSN");
        return new Response(JSON.stringify({ error: "No project ID found in the DSN" }), { status: 400 });
      }

      // Copy headers and set Host and Origin to the Sentry DSN
      const headers = new Headers(req.headers);
      headers.set("Host", dsn.host);
      headers.set("Origin", dsn.origin);

      // Tunnel the envelope to Sentry
      const upstream_sentry_url = `${dsn.origin}/api/${project_id}/envelope/`;

      log.info(`Tunneling to sentry ${dsn}`);
      log.info("Forwarding envelope", { sent_at: header.sent_at, event_id: header.event_id, sdk: header.sdk });

      // Send the envelope to Sentry
      return fetch(upstream_sentry_url, {
        method: "POST",
        headers: headers,
        body: envelopeBytes,
      })
    } catch (e) {
      log.error("Error tunneling to sentry", e);
      return new Response(JSON.stringify({ error: "Error tunneling to sentry" }), { status: 500 });
    }
  });
}
