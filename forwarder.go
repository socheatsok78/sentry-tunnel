package sentrytunnel

import "net/url"

func Forward(dsn *url.URL, envelope *Envelope) error {
	// fmt.Printf("Forwarding envelope: %v\n", envelope.Header.DSN)
	return nil
}
