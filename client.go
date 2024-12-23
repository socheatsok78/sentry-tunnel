package sentrytunnel

import (
	"net/http"
)

type SentryTunnelClient struct {
	Http      *http.Client
	UserAgent string
}

func (s *SentryTunnelClient) Forward(dsn SentryDSN, req http.Request) error {
	s.Http.Do(&req)
	return nil
}
