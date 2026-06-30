package incident

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
)

func TestNewOnBehalfOfClient(t *testing.T) {
	is := is.New(t)

	remoteHost := "https://test-stack.grafana.net/api/plugins/grafana-irm-app/resources/api/v1"
	accessToken := "test-access-token"
	idToken := "test-id-token"

	client := NewOnBehalfOfClient(remoteHost, accessToken, idToken)

	is.Equal(client.RemoteHost, remoteHost) // RemoteHost should be set
	is.True(client.HTTPClient != nil)       // HTTPClient should not be nil
	is.True(client.BeforeRequest != nil)    // BeforeRequest should not be nil
	is.True(client.Debug != nil)            // Debug should not be nil

	// Verify BeforeRequest sets the correct headers.
	req, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
	is.NoErr(err)

	err = client.BeforeRequest(req)
	is.NoErr(err)

	is.Equal(req.Header.Get("X-Access-Token"), accessToken) // X-Access-Token header
	is.Equal(req.Header.Get("X-Grafana-Id"), idToken)       // X-Grafana-Id header
	is.Equal(req.Header.Get("Authorization"), "")            // Authorization should not be set
}
