package sentrytunnel

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []byte
		expected    *Envelope
		expectedErr string
	}{
		{
			name:  "valid envelope",
			input: []byte(`{"dsn":"https://sentry.io/1"}` + "\n" + `{"type":"session"}` + "\n" + `{"message":"hello"}`),
			expected: &Envelope{
				Header: envelopeHeader{
					DSN: "https://sentry.io/1",
				},
				Type: envelopeMessageType{
					Type: "session",
				},
				Body: []byte(`{"message":"hello"}`),
			},
		},
		{
			name:        "invalid envelope",
			input:       []byte(`{"dsn":"https://sentry.io/1"}` + "\n" + `{"type":"session"}`),
			expected:    nil,
			expectedErr: "error parsing envelope",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			envelope, err := Parse(test.input)
			if err != nil {
				if test.expectedErr == "" {
					t.Errorf("unexpected error: %s", err)
				} else if err.Error() != test.expectedErr {
					t.Errorf("expected error: %s, got: %s", test.expectedErr, err)
				}
			}

			fmt.Println(envelope)

			if envelope != nil {
				if test.expected == nil {
					t.Errorf("expected nil, got: %v", envelope)
				} else if envelope.Header.DSN != test.expected.Header.DSN {
					t.Errorf("expected DSN: %s, got: %s", test.expected.Header.DSN, envelope.Header.DSN)
				} else if envelope.Type.Type != test.expected.Type.Type {
					t.Errorf("expected type: %s, got: %s", test.expected.Type.Type, envelope.Type.Type)
				} else if string(envelope.Body) != string(test.expected.Body) {
					t.Errorf("expected body: %s, got: %s", string(test.expected.Body), string(envelope.Body))
				}
			}
		})
	}
}
