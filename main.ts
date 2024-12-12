//    _____            __                ______                       __
//   / ___/___  ____  / /________  __   /_  __/_  ______  ____  ___  / /
//   \__ \/ _ \/ __ \/ __/ ___/ / / /    / / / / / / __ \/ __ \/ _ \/ / 
//  ___/ /  __/ / / / /_/ /  / /_/ /    / / / /_/ / / / / / / /  __/ /  
// /____/\___/_/ /_/\__/_/   \__, /    /_/  \__,_/_/ /_/_/ /_/\___/_/   
//                          /____/                                      
// A tunnel is an HTTP endpoint that acts as a proxy between Sentry and your application.
// Learn more at https://docs.sentry.io/platforms/javascript/troubleshooting/#using-the-tunnel-option

import * as log from "@std/log";

type SentryEnvelopeHeaderStruct = {
  dsn: string;
  sent_at: string;
  event_id: string;
  sdk: {
    name: string;
    version: string;
  };
}

type SentryEnvelopeTypeStruct = {
  type: string;
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
      const pieces = envelope.split("\n");
      const header: SentryEnvelopeHeaderStruct = JSON.parse(pieces[0]);
      const type: SentryEnvelopeTypeStruct = JSON.parse(pieces[1]);

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

      // Tunnel the envelope to Sentry
      const upstream_sentry_url = `${dsn.origin}/api/${project_id}/envelope/`;

      log.info("Forwarding envelope to sentry", {
        sent_at: header.sent_at,
        event_id: header.event_id,
        type: type.type,
        sdk: header.sdk,
        envelope_byte_length: envelopeBytes.byteLength,
      });

      // Send the envelope to Sentry
      return fetch(upstream_sentry_url, {
        method: "POST",
        body: envelopeBytes,
      })
    } catch (e) {
      log.error("Error tunneling to sentry", e);
      return new Response(JSON.stringify({ error: "Error tunneling to sentry" }), { status: 500 });
    }
  });
}
