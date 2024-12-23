package sentrytunnel

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type SentryDSN struct {
	*url.URL
}

type Envelope struct {
	Header envelopeHeader
	Type   envelopeMessageType
	Body   envelopeMessageBody
	Data   []byte
}

type sdk struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Parse(bytes []byte) (*Envelope, error) {
	envelope := &Envelope{}
	err := Unmarshal(bytes, envelope)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}

func Unmarshal(bytes []byte, envelope *Envelope) error {
	lines := strings.SplitN(string(bytes), "\n", 3)

	// Ensure the envelope has at least 3 lines
	// The first line is the header, the second line is the type, and the third line is the body
	if len(lines) < 3 {
		return fmt.Errorf("error parsing envelope")
	}

	// Parse the envelope header
	envelopeHeader, err := parseEnvelopeHeader([]byte(lines[0]))
	if err != nil {
		return err
	}

	// Parse the envelope type
	envelopeType, err := parseEnvelopeType([]byte(lines[1]))
	if err != nil {
		return err
	}

	envelope.Header = *envelopeHeader
	envelope.Type = envelopeType
	envelope.Body = envelopeMessageBody(lines[2])
	envelope.Data = bytes

	return nil
}

// Parse parses the given bytes into an EnvelopeHeader
func parseEnvelopeHeader(bytes []byte) (*envelopeHeader, error) {
	envelopeHeader := &envelopeHeader{}

	err := json.Unmarshal(bytes, envelopeHeader)
	if err != nil {
		return nil, fmt.Errorf("error parsing envelope header")
	}

	return envelopeHeader, nil
}

// Parse parses the given bytes into an EnvelopeMessageType
func parseEnvelopeType(bytes []byte) (envelopeMessageType, error) {
	envelopeMessageType := envelopeMessageType{}

	err := json.Unmarshal(bytes, &envelopeMessageType)
	if err != nil {
		return envelopeMessageType, fmt.Errorf("error parsing envelope type")
	}

	return envelopeMessageType, nil
}

type envelopeHeader struct {
	DSN     string `json:"dsn"`
	EventID string `json:"event_id"`
	SentAt  string `json:"sent_at"`
	SDK     sdk    `json:"sdk"`
}

type envelopeMessageType struct {
	Type string `json:"type"`
}

type envelopeMessageBody []byte
