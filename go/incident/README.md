# Grafana Incident Go Client library

The Grafana Incident Go Client library allows you to access the Grafana Incident API from your Go code.

- Get started with the [Grafana Incident API documentation (preview)](https://grafana.com/docs/grafana-cloud/incident/api/experimental/)
- Or dive deep into the [reference docs (preview)](https://grafana.com/docs/grafana-cloud/incident/api/experimental/reference/)

## Get started

Import the package:

```
go get github.com/grafana/incident-api/go/incident
```

## Use the API to create an Incident

```go
// create a client, and the services you need
serviceAccountToken := os.Getenv("SERVICE_ACCOUNT_TOKEN")
client := incident.NewClient("https://your-api-endpoint/api", serviceAccountToken)
incidentsService := incident.NewIncidentsService(client)

// declare an incident
createIncidentResp, err := incidentsService.CreateIncident(ctx, incident.CreateIncidentRequest{
	Title: "short description explaining what's going wrong",
})
if err != nil {
	// if something goes wrong, the error will help you
	return fmt.Errorf("create incident: %w", err)
}
// success, get the details from the createIncidentResp object
fmt.Println("declared Incident", createIncidentResp.Incident.IncidentID)
```

## Write an Outgoing Webhook handler

You can use the [Outgoing Webhook integration](https://grafana.com/docs/grafana-cloud/incident/integrations/configure-outgoing-webhooks/) to get Grafana Incident to POST a request on specific events.

If you are consuming that event in Go, you can use the `incident.ParseWebhook` helper:

```go
import (
	"github.com/grafana/incident-api/go/incident"
)

// handleIncidentWebhook gets a handler that processes webhooks from
// Grafana Incident.
func handleIncidentWebhook(secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// incident.ParseWebhook will verify the signature and decode
		// the body into the incident.OutgoingWebhookPayload type.
		payload, err := incident.ParseWebhook(r, secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch payload.Event {
		case "grafana.incident.created":
			fmt.Printf("Incident declared: %s\n", payload.Incident.Title)
		default:
			fmt.Printf("Unknown event: %s\n", payload.Event)
		}
	}))
}
```
