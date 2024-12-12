//    _____            __                ______                       __
//   / ___/___  ____  / /________  __   /_  __/_  ______  ____  ___  / /
//   \__ \/ _ \/ __ \/ __/ ___/ / / /    / / / / / / __ \/ __ \/ _ \/ / 
//  ___/ /  __/ / / / /_/ /  / /_/ /    / / / /_/ / / / / / / /  __/ /  
// /____/\___/_/ /_/\__/_/   \__, /    /_/  \__,_/_/ /_/_/ /_/\___/_/   
//                          /____/                                      
// A tunnel is an HTTP endpoint that acts as a proxy between Sentry and your application.
// Learn more at https://docs.sentry.io/platforms/javascript/troubleshooting/#using-the-tunnel-option

if (import.meta.main) {
  const options: Deno.ServeTcpOptions = {
    port: parseInt(Deno.env.get("PORT") as string || "3003"),
  }

  Deno.serve(options, async (req) => {
    if (req.method !== "POST") {
      return new Response(null, { status: 405 });
    }
    try {
      const envelopeBytes = await req.arrayBuffer();
      const envelope = new TextDecoder().decode(envelopeBytes);
      const piece = envelope.split("\n")[0];
      const header = JSON.parse(piece);

      if (!header["dsn"]) {
        console.error("No DSN found in the envelope header");
        return new Response(JSON.stringify({ error: "No DSN found in the envelope header" }), { status: 400 });
      }

      const dsn = new URL(header["dsn"]);
      const project_id = dsn.pathname?.replace("/", "");

      if (!project_id) {
        console.error("No project ID found in the DSN");
        return new Response(JSON.stringify({ error: "No project ID found in the DSN" }), { status: 400 });
      }

      console.info("Tunneling to sentry", `${dsn}`, );

      const upstream_sentry_url = `${dsn.origin}/api/${project_id}/envelope/`;
      return fetch(upstream_sentry_url, {
        method: "POST",
        body: envelopeBytes,
      });
    } catch (e) {
      console.error("Error tunneling to sentry", e);
      return new Response(JSON.stringify({ error: "error tunneling to sentry" }), { status: 500 });
    }
  });
}
