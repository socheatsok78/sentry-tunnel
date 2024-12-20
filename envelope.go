package sentrytunnel

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Envelope struct {
	Header envelopeHeader
	Type   envelopeMessageType
	Body   envelopeMessageBody
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

type sdk struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Parse(bytes []byte) (*Envelope, error) {
	envelopeStr := string(bytes)
	lines := strings.SplitN(envelopeStr, "\n", 3)

	// Ensure the envelope has at least 3 lines
	// The first line is the header, the second line is the type, and the third line is the body
	if len(lines) < 3 {
		return nil, fmt.Errorf("error parsing envelope")
	}

	// Parse the envelope header
	envelopeHeader, err := parseEnvelopeHeader([]byte(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("error parsing envelope header")
	}

	// Parse the envelope type
	envelopeType, err := parseEnvelopeType([]byte(lines[1]))
	if err != nil {
		return nil, fmt.Errorf("error parsing envelope type")
	}

	// Create the envelope object
	envelope := Envelope{
		Header: *envelopeHeader,
		Type:   envelopeType,
		Body:   envelopeMessageBody(lines[2]),
	}

	return &envelope, nil
}

// Parse parses the given bytes into an EnvelopeHeader
func parseEnvelopeHeader(bytes []byte) (*envelopeHeader, error) {
	EnvelopeHeader := &envelopeHeader{}

	err := json.Unmarshal(bytes, EnvelopeHeader)
	if err != nil {
		return nil, fmt.Errorf("error parsing envelope header: %w", err)
	}

	return EnvelopeHeader, nil
}

// Parse parses the given bytes into an EnvelopeMessageType
func parseEnvelopeType(bytes []byte) (envelopeMessageType, error) {
	EnvelopeMessageType := envelopeMessageType{}

	err := json.Unmarshal(bytes, &EnvelopeMessageType)
	if err != nil {
		return EnvelopeMessageType, fmt.Errorf("error parsing envelope type: %w", err)
	}

	return EnvelopeMessageType, nil
}
