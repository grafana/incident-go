package incident

import (
	"net/http"
	"time"
)

// NewOnBehalfOfClient makes a new Client that authenticates on behalf of
// a specific Grafana user using their access token and ID token.
//
// This is used when a service needs to make API calls as a specific user,
// rather than as a service account. The accessToken is set as the
// X-Access-Token header and the idToken is set as the X-Grafana-Id header
// on every request.
//
// The remoteHost should be
// "https://your-stack.grafana.net/api/plugins/grafana-irm-app/resources/api/v1"
// with `your-stack.grafana.net` pointing to your instance.
func NewOnBehalfOfClient(remoteHost, accessToken, idToken string) *Client {
	c := &Client{
		RemoteHost: remoteHost,
		Debug:      func(s string) { /* no-op by default */ },
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		BeforeRequest: func(r *http.Request) error {
			r.Header.Set("X-Access-Token", accessToken)
			r.Header.Set("X-Grafana-Id", idToken)
			return nil
		},
	}
	return c
}
