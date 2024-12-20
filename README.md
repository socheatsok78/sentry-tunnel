![Sentry](https://camo.githubusercontent.com/ebaaa2f1ab4b8efc7284c9f736e0de6f8fca212a6a14e0255ad1706d1e80f76d/68747470733a2f2f73656e7472792d6272616e642e73746f726167652e676f6f676c65617069732e636f6d2f73656e7472792d776f72646d61726b2d6461726b2d3238307838342e706e67)

> [!NOTE]
> The implementation is based on the explanation provided by the [official sentry documentation][sentry-tunnel-docs].

## Sentry Tunnel

A tunnel is an HTTP endpoint that acts as a proxy between Sentry and your application.

Because you control this server, there is no risk of any requests sent to it being blocked. When the endpoint lives under the same origin (although it does not have to in order for the tunnel to work), the browser will not treat any requests to the endpoint as a third-party request. As a result, these requests will have different security measures applied which, by default, don't trigger ad-blockers.

A quick summary of the flow can be found below.

![tunnel.png](https://docs.sentry.io/_next/image/?url=%2Fmdx-images%2Ftunnel-7ZZLHFR5.png%231374x1078&w=1920&q=75)

Starting with version `6.7.0` of the JavaScript SDK, you can use the `tunnel` option to tell the SDK to deliver events to the configured URL, instead of using the DSN. This allows the SDK to remove `sentry_key` from the query parameters, which is one of the main reasons ad-blockers prevent sending events in the first place. This option also stops the SDK from sending preflight requests, which was one of the requirements that necessitated sending the `sentry_key` in the query parameters.

To enable the `tunnel` option, provide either a relative or an absolute URL in your `Sentry.init` call. When you use a relative URL, it's relative to the current origin, and this is the form that we recommend. Using a relative URL will not trigger a preflight CORS request, so no events will be blocked, because the ad-blocker will not treat these events as third-party requests.

```js
Sentry.init({
  dsn: "https://examplePublicKey@o0.ingest.sentry.io/0",
  tunnel: "/tunnel",
});
```

Once configured, all events will be sent to the /tunnel endpoint. This solution, however, requires an additional configuration on the server, as the events now need to be parsed and redirected to Sentry.

## Metrics
The tunnel server exposes the following metrics:

```
# HELP sentry_envelope_accepted The number of envelopes accepted by the tunnel
# TYPE sentry_envelope_accepted counter
sentry_envelope_accepted 0
# HELP sentry_envelope_forwarded The number of envelopes forwarded by the tunnel
# TYPE sentry_envelope_forwarded counter
sentry_envelope_forwarded 0
# HELP sentry_envelope_forwarded_error The number of envelopes that failed to be forwarded by the tunnel
# TYPE sentry_envelope_forwarded_error counter
sentry_envelope_forwarded_error 0
# HELP sentry_envelope_rejected The number of envelopes rejected by the tunnel
# TYPE sentry_envelope_rejected counter
sentry_envelope_rejected 0
```

## Benchmark
The benchmark was done using [wrk](https://github.com/wg/wrk) running on a local machine.

```
Running 30s test @ http://localhost:8080/tunnel
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.74ms    5.41ms 100.42ms   86.42%
    Req/Sec     6.51k     1.24k   14.61k    77.06%
  2337038 requests in 30.08s, 167.16MB read
  Socket errors: connect 0, read 376, write 0, timeout 0
Requests/sec:  77685.87
Transfer/sec:      5.56MB
```

<!-- Links -->
[sentry-tunnel-docs]: https://docs.sentry.io/platforms/javascript/troubleshooting/#using-the-tunnel-option
