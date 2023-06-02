package incident

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ParseWebhook parses an Outgoing Webhook from Grafana Incident.
// Use this function when writing code to handle the event.
// The signature will be verified using the known secret.
func ParseWebhook(r *http.Request, signingSecret string) (*OutgoingWebhookPayload, error) {
	err := VerifySignature(r, signingSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %w", err)
	}
	payload, err := io.ReadAll(io.LimitReader(r.Body, int64(1*mb)))
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	var payloadObject OutgoingWebhookPayload
	err = json.Unmarshal(payload, &payloadObject)
	if err != nil {
		return nil, fmt.Errorf("unmarshal webhook payload: %w", err)
	}
	return &payloadObject, nil
}

// VerifySignature checks Gi-Signature header against the secret you
// got when enabling the integration in the tool.
func VerifySignature(r *http.Request, signingSecret string) error {
	header := r.Header["Gi-Signature"]
	if len(header) == 0 || header[0] == "" {
		return errors.New("empty GI-Signature")
	}
	signatures := strings.Split(header[0], ",")
	s := make(map[string]string)
	for _, pair := range signatures {
		values := strings.Split(pair, "=")
		s[values[0]] = values[1]
	}
	t := s["t"]
	v1 := s["v1"]
	payload, err := io.ReadAll(io.LimitReader(r.Body, int64(1*mb)))
	// Copy body for other handlers
	r.Body = io.NopCloser(bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	bodyHash := Hash(payload)
	stringToSign := bodyHash + ":" + t + ":v1"
	expected := GenerateSignature([]byte(stringToSign), signingSecret)
	if expected != v1 {
		return errors.New("invalid GI-Signature")
	}
	return nil
}

// Hash encodes SHA256 hash to Base64.
func Hash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// GenerateSignature creates SHA256 hash.
func GenerateSignature(data []byte, secret string) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// Byte size size suffixes.
const (
	// mb represents megabytes.
	mb = 1 << (10 * 2)
)
