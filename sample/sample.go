package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func postAction(request *http.Request) (int, []byte) {
	envelopeBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	envelope := string(envelopeBytes)
	lines := strings.SplitN(envelope, "\n", 2)
	if len(lines) < 1 {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	header := make(map[string]interface{})
	if err := json.Unmarshal([]byte(lines[0]), &header); err != nil {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	dsn, ok := header["dsn"].(string)
	if !ok {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	fmt.Printf("Received envelope for DSN: %s\n", dsn)

	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	projectID := strings.TrimPrefix(dsnURL.Path, "/")

	upstreamSentryURL := "https://" + dsnURL.Host + "/api/" + projectID + "/envelope/"
	resp, err := http.Post(upstreamSentryURL, "application/octet-stream", bytes.NewReader(envelopeBytes))
	if err != nil || resp.StatusCode != 200 {
		return 500, []byte(`{"error":"error tunneling to sentry"}`)
	}

	return 200, []byte("{}")
}

func main() {
	http.HandleFunc("/tunnel", func(w http.ResponseWriter, r *http.Request) {
		status, response := postAction(r)
		w.WriteHeader(status)
		w.Write(response)
	})

	http.ListenAndServe(":8080", nil)
}
