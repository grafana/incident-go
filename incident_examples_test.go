package incident_test

import (
	"context"
	"fmt"
	"log"

	"github.com/grafana/incident-api/go/incident"
)

// ExampleCreateIncident shows how to create an incident.
// It uses the incident.NewTestClient so all responses are stubbed,
// you should use incident.NewClient, specifying the API endpoint
// and the service account token to send with the requests.
func ExampleCreateIncident() {
	ctx := context.Background()
	client := incident.NewTestClient()
	incidentsService := incident.NewIncidentsService(client)
	createIncidentResp, err := incidentsService.CreateIncident(ctx, incident.CreateIncidentRequest{
		Title:    "high latency in web requests",
		Severity: incident.Options.IncidentSeverity.Minor,
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("new incident: %s\n", createIncidentResp.Incident.Title)
	// Output: new incident: high latency in web requests
}
