package incident

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestVerifySignature(t *testing.T) {
	is := is.New(t)
	method := "POST"
	url := "https://some-test-url"
	secret := "some-secret-value"

	goodRequest, err := http.NewRequest(method, url, strings.NewReader(testPayload))
	is.NoErr(err)
	goodRequest.Header.Set("gi-signature", "t=1678360185,v1=66ed620915831e9ab71604bf35d9b7cf2e71978ed6e62be6fd0e5194289849c7")
	err = VerifySignature(goodRequest, secret) // valid gi-signature
	is.NoErr(err)

	badRequest, err := http.NewRequest(method, url, strings.NewReader(testPayload))
	is.NoErr(err)
	badRequest.Header.Set("gi-signature", "t=123456,v1=some-bad-signature")
	err = VerifySignature(badRequest, secret)
	is.Equal(err, errors.New("invalid GI-Signature")) // invalid gi-signature
}

func TestParseWebhook(t *testing.T) {
	is := is.New(t)
	testURL := "https://some-test-url"
	testSecret := "some-secret-value"
	timestamp, err := time.Parse(time.RFC3339, "2023-03-09T11:09:45Z")
	unixTime := fmt.Sprintf("%d", timestamp.Unix())
	is.NoErr(err)

	request, err := http.NewRequest("POST", testURL, strings.NewReader(testPayload))
	is.NoErr(err)
	bodyHash := Hash([]byte(testPayload))
	stringToSign := bodyHash + ":" + unixTime + ":v1"
	signature := GenerateSignature([]byte(stringToSign), testSecret)
	request.Header.Set("GI-Signature", fmt.Sprintf("t=%s,v1=%s", unixTime, signature))
	parsed, err := ParseWebhook(request, testSecret)
	is.NoErr(err)
	is.Equal(parsed.ID, "webhook-out-6409be79a7ee628da2a70dbe") // parsed.ID
}

var testPayload = `{
	"version": "v1.0.0",
	"id": "webhook-out-6409be79a7ee628da2a70dbe",
	"source": "/grafana/incident",
	"time": "2023-03-09T11:09:45Z",
	"event": "grafana.incident.updated.role",
	"incident": {
		"incidentID": "19",
		"severity": "pending",
		"labels": [],
		"isDrill": false,
		"createdTime": "2023-03-09T11:09:44.293367Z",
		"modifiedTime": "2023-03-09T11:09:45.602445Z",
		"createdByUser": {
		  "userID": "grafana-incident:user-637c92ccd784a2ee",
		  "name": "paulcoghlan",
		  "photoURL": "https://www.gravatar.com/avatar/5bf54eaa42d21c0cd36ec62ac79f3e28?s=512&d=https%3A%2F%2Favatars.slack-edge.com%2F2022-09-26%2F4130143663477_91b972302da73fcbbd21_192.png"
		},
		"closedTime": "",
		"durationSeconds": 1,
		"status": "active",
		"title": "1111",
		"overviewURL": "/a/grafana-incident-app/incidents/19/1111",
		"roles": [
		  {
			"role": "observer",
			"description": "Watching the incident",
			"maxPeople": 0,
			"mandatory": false,
			"important": false,
			"user": {
			  "userID": "grafana-incident:user-637c92ccd784a2ee",
			  "name": "paulcoghlan",
			  "photoURL": "https://www.gravatar.com/avatar/5bf54eaa42d21c0cd36ec62ac79f3e28?s=512&d=https%3A%2F%2Favatars.slack-edge.com%2F2022-09-26%2F4130143663477_91b972302da73fcbbd21_192.png"
			}
		  }
		],
		"taskList": {
		  "tasks": [
			{
			  "taskID": "assign-role-investigator",
			  "immutable": true,
			  "createdTime": "2023-03-09T11:09:44.216113029Z",
			  "modifiedTime": "2023-03-09T11:09:44.216124254Z",
			  "text": "Assign INVESTIGATOR role",
			  "status": "todo",
			  "authorUser": null,
			  "assignedUser": null
			},
			{
			  "taskID": "assign-role-commander",
			  "immutable": true,
			  "createdTime": "2023-03-09T11:09:44.216128351Z",
			  "modifiedTime": "2023-03-09T11:09:44.216129374Z",
			  "text": "Assign COMMANDER role",
			  "status": "todo",
			  "authorUser": null,
			  "assignedUser": null
			},
			{
			  "taskID": "must-choose-severity",
			  "immutable": true,
			  "createdTime": "2023-03-09T11:09:44.216131603Z",
			  "modifiedTime": "2023-03-09T11:09:44.216132579Z",
			  "text": "Specify incident severity",
			  "status": "todo",
			  "authorUser": null,
			  "assignedUser": null
			}
		  ],
		  "todoCount": 3,
		  "doneCount": 0
		},
		"summary": "",
		"heroImagePath": "/api/hero-images/13/FcmfxKrsFo4hY55dqi1zLJLM4r13On3XyUmeQEb5yFhhq4UzjLWkQ9YgOjUO5uVmgyS9AEJNzyx9bzY3iWNkBeMSAQuyVu1U3NdMrGkMoVlujgj4u5ibpg9eV2CHLX8C/v10/19.png",
		"incidentStart": "2023-03-09T11:09:44.293367Z",
		"incidentEnd": ""
	  }
	}`
