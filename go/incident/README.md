# Grafana Incident Go Client library

The Grafana Incident Go Client library allows you to access the Grafana Incident API from your Go code.

- Get started with the [Grafana Incident API documentation (preview)](https://grafana.com/docs/grafana-cloud/incident/api/preview/)
- Or dive deep into the [reference docs (preview)](https://grafana.com/docs/grafana-cloud/incident/api/preview/reference/)

## Get started

Import the package:

```
go get github.com/grafana/incident-api/go/incident
```

Simple example:

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
