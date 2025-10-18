package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/models"
)

// CreateScheduleRequest defines the request body for creating a new schedule.
type CreateScheduleRequest struct {
	Name           string `json:"name" validate:"required,min=1"`
	Minute         string `json:"minute" validate:"required"`
	Hour           string `json:"hour" validate:"required"`
	DayOfMonth     string `json:"day_of_month" validate:"required"`
	Month          string `json:"month" validate:"required"`
	DayOfWeek      string `json:"day_of_week" validate:"required"`
	IsActive       bool   `json:"is_active,omitempty"`
	OnlyWhenOnline bool   `json:"only_when_online,omitempty"`
}

// UpdateScheduleRequest defines the request body for updating a schedule.
type UpdateScheduleRequest struct {
	Name           string `json:"name,omitempty" validate:"omitempty,min=1"`
	Minute         string `json:"minute,omitempty"`
	Hour           string `json:"hour,omitempty"`
	DayOfMonth     string `json:"day_of_month,omitempty"`
	Month          string `json:"month,omitempty"`
	DayOfWeek      string `json:"day_of_week,omitempty"`
	IsActive       *bool  `json:"is_active,omitempty"`
	OnlyWhenOnline *bool  `json:"only_when_online,omitempty"`
}

// ListSchedules retrieves all scheduled tasks for a server.
func (c *client) ListSchedules(ctx context.Context, serverID string) ([]*models.Schedule, error) {
	path := fmt.Sprintf("client/servers/%s/schedules", serverID)
	var response struct {
		Data []struct {
			Attributes models.Schedule `json:"attributes"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}

	schedules := make([]*models.Schedule, len(response.Data))
	for i, item := range response.Data {
		schedules[i] = &item.Attributes
	}
	return schedules, nil
}

// GetSchedule retrieves details for a specific schedule.
func (c *client) GetSchedule(ctx context.Context, serverID string, scheduleID int) (*models.Schedule, error) {
	path := fmt.Sprintf("client/servers/%s/schedules/%d", serverID, scheduleID)
	var response struct {
		Attributes models.Schedule `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// CreateSchedule creates a new scheduled task.
func (c *client) CreateSchedule(ctx context.Context, serverID string, req CreateScheduleRequest) (*models.Schedule, error) {
	path := fmt.Sprintf("client/servers/%s/schedules", serverID)
	var response struct {
		Attributes models.Schedule `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// UpdateSchedule modifies an existing schedule.
func (c *client) UpdateSchedule(ctx context.Context, serverID string, scheduleID int, req UpdateScheduleRequest) (*models.Schedule, error) {
	path := fmt.Sprintf("client/servers/%s/schedules/%d", serverID, scheduleID)
	var response struct {
		Attributes models.Schedule `json:"attributes"`
	}
	_, err := c.client.Do(ctx, http.MethodPost, path, req, &response)
	if err != nil {
		return nil, err
	}
	return &response.Attributes, nil
}

// DeleteSchedule permanently deletes a schedule and its tasks.
func (c *client) DeleteSchedule(ctx context.Context, serverID string, scheduleID int) error {
	path := fmt.Sprintf("client/servers/%s/schedules/%d", serverID, scheduleID)
	_, err := c.client.Do(ctx, http.MethodDelete, path, nil, nil)
	return err
}

// ExecuteSchedule manually triggers a schedule to run immediately.
func (c *client) ExecuteSchedule(ctx context.Context, serverID string, scheduleID int) error {
	path := fmt.Sprintf("client/servers/%s/schedules/%d/execute", serverID, scheduleID)
	_, err := c.client.Do(ctx, http.MethodPost, path, nil, nil)
	return err
}
