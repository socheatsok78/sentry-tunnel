package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sentrytunnel "github.com/socheatsok78/sentry-tunnel"
	"github.com/urfave/cli/v3"
)

var (
	logger log.Logger
)

var (
	// SentryEnvelopeAccepted is a Prometheus counter for the number of envelopes accepted by the tunnel
	SentryEnvelopeAccepted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentry_envelope_accepted",
		Help: "The number of envelopes accepted by the tunnel",
	})
	// SentryEnvelopeRejected is a Prometheus counter for the number of envelopes rejected by the tunnel
	SentryEnvelopeRejected = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentry_envelope_rejected",
		Help: "The number of envelopes rejected by the tunnel",
	})
	// SentryEnvelopeForwardedSuccess is a Prometheus counter for the number of envelopes successfully forwarded by the tunnel
	SentryEnvelopeForwardedSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentry_envelope_forward_success",
		Help: "The number of envelopes successfully forwarded by the tunnel",
	})
	// SentryEnvelopeForwardedError is a Prometheus counter for the number of envelopes that failed to be forwarded by the tunnel
	SentryEnvelopeForwardedError = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentry_envelope_forward_error",
		Help: "The number of envelopes that failed to be forwarded by the tunnel",
	})
)

func init() {
	// Set up logging
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Register Prometheus metrics
	prometheus.MustRegister(SentryEnvelopeAccepted)
	prometheus.MustRegister(SentryEnvelopeRejected)
	prometheus.MustRegister(SentryEnvelopeForwardedSuccess)
	prometheus.MustRegister(SentryEnvelopeForwardedError)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := cli.Command{
		Name:  "sentry-tunnel",
		Usage: "A tunneling service for Sentry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen-addr",
				Usage: "The address to listen on",
				Value: ":8080",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Set the log level",
				Value: "info",
			},
			&cli.StringSliceFlag{
				Name:  "trusted-sentry-dsn",
				Usage: "A map of Sentry DSNs that are trusted by the tunnel",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error { return action(ctx, c) },
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		panic(err)
	}
}

func action(_ context.Context, cmd *cli.Command) error {
	listenAddr := cmd.String("listen-addr")
	trustedDSNs := cmd.StringSlice("trusted-sentry-dsn")

	// Register Prometheus metrics handler
	http.Handle("/metrics", promhttp.Handler())

	// Register the tunnel handler
	http.Handle("/tunnel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		envelopeBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"error tunneling to sentry"}`))
			return
		}

		envelope, err := sentrytunnel.Parse(envelopeBytes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"error tunneling to sentry"}`))
			return
		}

		for _, trustedDSN := range trustedDSNs {
			if envelope.Header.DSN == trustedDSN {
				break
			} else {
				SentryEnvelopeRejected.Inc()
				w.WriteHeader(403)
				w.Write([]byte(`{"error":"untrusted DSN"}`))
				return
			}
		}

		SentryEnvelopeAccepted.Inc()

		if err := sentrytunnel.Forward(envelope); err != nil {
			SentryEnvelopeForwardedError.Inc()
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"error tunneling to sentry"}`))
			return
		}

		SentryEnvelopeForwardedSuccess.Inc()
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	// Start the server
	level.Info(logger).Log("msg", "The tunnel is now listening on: "+listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}
