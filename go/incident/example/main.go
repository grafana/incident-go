package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grafana/incident-api/go/incident"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("missing secret")
		os.Exit(1)
	}
	secret := os.Args[1]
	auth := Parse(secret)
	s := &http.Server{
		Addr:    ":3001",
		Handler: auth,
	}
	log.Fatal(s.ListenAndServe())
}

// Parse logs incident title using webhook
func Parse(secret string) http.Handler {
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
	})
}
