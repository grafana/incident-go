# Grafana Incident API - Go client library

The Grafana Incident Go client library allows you to access the Grafana Incident API from your Go code.

- Get started with the [Grafana Incident API documentation](https://grafana.com/docs/grafana-cloud/incident/api/)
- Or dive deep into the [Go reference docs](https://grafana.com/docs/grafana-cloud/incident/api/reference/go/)

## Get started

Import the package:

```
go get github.com/grafana/incident-go@latest
```

## Make calls to the API

In this example, we will use the [IncidentsService.CreateIncident() method](https://grafana.com/docs/grafana-cloud/incident/api/experimental/reference/go/#createincident) to declare an Incident, and print its ID.

```go
// create a client, and the services you need
serviceAccountToken := os.Getenv("SERVICE_ACCOUNT_TOKEN")
client := incident.NewClient("https://your-api-endpoint/api", serviceAccountToken)
incidentsService := incident.NewIncidentsService(client)

// declare an incident
createIncidentResp, err := incidentsService.CreateIncident(ctx, incident.CreateIncidentRequest{
	Title: "short description explaining what's going wrong",
	Severity: incident.Options.IncidentSeverity.Minor,
})
if err != nil {
	// if something goes wrong, the error will help you
	return fmt.Errorf("create incident: %w", err)
}
// success, get the details from the createIncidentResp object
fmt.Println("declared Incident", createIncidentResp.Incident.IncidentID)
```

## Handle webhooks from Grafana Incident

You can use the [Outgoing Webhook integration](https://grafana.com/docs/grafana-cloud/incident/integrations/configure-outgoing-webhooks/) to get Grafana Incident to POST a request on specific events.

If you are consuming that event in Go, you can use the `incident.ParseWebhook` helper:

```go
import (
	incident "github.com/grafana/incident-go"
)

// handleIncidentWebhook gets a handler that processes webhooks from
// Grafana Incident.
// The secret should be safely injected (avoid committing it to source control).
// Secrets can be created in the web interface when configuring the Outgoing Webhook integration.
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
