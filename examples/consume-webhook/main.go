package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	incident "github.com/grafana/incident-go"
)

func main() {
	secret := os.Getenv("GRAFANA_INCIDENT_SIGNING_SECRET")
	s := &http.Server{
		Addr:    ":3001",
		Handler: incidentWebhookHandler(secret),
	}
	log.Fatal(s.ListenAndServe())
}

// incidentWebhookHandler logs the title of the incident
// by parsing the webhook.
func incidentWebhookHandler(secret string) http.Handler {
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
			fmt.Printf("Ignoring event: %s\n", payload.Event)
		}
	})
}
