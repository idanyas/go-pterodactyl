package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/idanyas/go-pterodactyl/client"
)

func TestSchedules_ListSchedules(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [{
				"object": "server_schedule",
				"attributes": {
					"id": 1,
					"name": "Daily Restart",
					"cron": {
						"minute": "0",
						"hour": "3",
						"day_of_month": "*",
						"month": "*",
						"day_of_week": "*"
					},
					"is_active": true
				}
			}]
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	schedules, err := clientAPI.ListSchedules(context.Background(), "d3aac109")
	if err != nil {
		t.Fatalf("ListSchedules() error = %v", err)
	}

	if len(schedules) != 1 {
		t.Errorf("expected 1 schedule, got %d", len(schedules))
	}
}

func TestSchedules_CreateSchedule(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	req := client.CreateScheduleRequest{
		Name:       "Hourly Backup",
		Minute:     "0",
		Hour:       "*",
		DayOfMonth: "*",
		Month:      "*",
		DayOfWeek:  "*",
	}

	mux.HandleFunc("/api/client/servers/d3aac109/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var received client.CreateScheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}

		if received.Name != req.Name {
			t.Errorf("Name = %s, want %s", received.Name, req.Name)
		}

		fmt.Fprint(w, `{
			"object": "server_schedule",
			"attributes": {
				"id": 2,
				"name": "Hourly Backup"
			}
		}`)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	schedule, err := clientAPI.CreateSchedule(context.Background(), "d3aac109", req)
	if err != nil {
		t.Fatalf("CreateSchedule() error = %v", err)
	}

	if schedule.Name != req.Name {
		t.Errorf("Name = %s, want %s", schedule.Name, req.Name)
	}
}

func TestSchedules_DeleteSchedule(t *testing.T) {
	mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/client/servers/d3aac109/schedules/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	c := testClient(t, serverURL)
	clientAPI := client.New(c)

	err := clientAPI.DeleteSchedule(context.Background(), "d3aac109", 1)
	if err != nil {
		t.Fatalf("DeleteSchedule() error = %v", err)
	}
}
