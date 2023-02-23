package incident_test

import (
	"context"
	"testing"

	"github.com/grafana/incident-api/golang/incident"
	"github.com/matryer/is"
)

func Test(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	client := incident.NewTestClient()
	incidentsService := incident.NewIncidentsService(client)
	createIncidentResp, err := incidentsService.CreateIncident(ctx, incident.CreateIncidentRequest{
		Title: "short description explaining what's going wrong",
	})
	is.NoErr(err)
	is.Equal(createIncidentResp.Incident.IncidentID, "incident-123")

}
