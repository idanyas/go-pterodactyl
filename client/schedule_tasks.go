package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// CreateScheduleTaskRequest defines the request body for creating a schedule task.
type CreateScheduleTaskRequest struct {
	Action            string `json:"action" validate:"required,oneof=command power backup"`
	Payload           string `json:"payload" validate:"required"`
	TimeOffset        int    `json:"time_offset" validate:"gte=0"`
	ContinueOnFailure bool   `json:"continue_on_failure,omitempty"`
}

// UpdateScheduleTaskRequest defines the request body for updating a schedule task.
type UpdateScheduleTaskRequest struct {
	Action            string `json:"action,omitempty" validate:"omitempty,oneof=command power backup"`
	Payload           string `json:"payload,omitempty"`
	TimeOffset        int    `json:"time_offset,omitempty" validate:"omitempty,gte=0"`
	ContinueOnFailure *bool  `json:"continue_on_failure,omitempty"`
}

// CreateScheduleTask adds a new task to a schedule.
func (c *client) CreateScheduleTask(ctx context.Context, serverID string, scheduleID int, req CreateScheduleTaskRequest) (*models.ScheduleTask, error) {
	path := fmt.Sprintf("client/servers/%s/schedules/%d/tasks", serverID, scheduleID)
	var response struct {
		Attributes models.ScheduleTask `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// UpdateScheduleTask modifies an existing task within a schedule.
func (c *client) UpdateScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int, req UpdateScheduleTaskRequest) (*models.ScheduleTask, error) {
	path := fmt.Sprintf("client/servers/%s/schedules/%d/tasks/%d", serverID, scheduleID, taskID)
	var response struct {
		Attributes models.ScheduleTask `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DeleteScheduleTask removes a task from a schedule.
func (c *client) DeleteScheduleTask(ctx context.Context, serverID string, scheduleID, taskID int) error {
	path := fmt.Sprintf("client/servers/%s/schedules/%d/tasks/%d", serverID, scheduleID, taskID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
