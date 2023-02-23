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

func TestAddTask(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	client := incident.NewTestClient()
	tasksService := incident.NewTasksService(client)
	addTaskResp, err := tasksService.AddTask(ctx, incident.AddTaskRequest{
		IncidentID: "incident-123",
		Text:       "Consider tweeting to let people know",
	})
	is.NoErr(err)
	is.Equal(addTaskResp.Task.TaskID, "task-123456")

}

func TestCheckLabels(t *testing.T) {

	is := is.New(t)

	ctx := context.Background()
	client := incident.NewTestClient()
	incidentsService := incident.NewIncidentsService(client)
	getIncidentResp, err := incidentsService.GetIncident(ctx, incident.GetIncidentRequest{
		IncidentID: "abc123",
	})
	is.NoErr(err)
	for _, item := range getIncidentResp.Incident.Labels {
		if item.Label == "specific-label" {
			postSlackMessage(ctx, "#security-channel", "An incident has been flagged with the security label.")
			return
		}
	}

}

func postSlackMessage(context.Context, string, string) {}
