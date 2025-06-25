package api

import "time"

type Task struct {
	ID                int       `json:"id"`
	SequenceID        int       `json:"sequence_id"`
	Action            string    `json:"action"` // e.g., "command", "power", "backup"
	Payload           string    `json:"payload"`
	TimeOffset        int       `json:"time_offset"` // Seconds to wait after the previous task
	IsQueued          bool      `json:"is_queued"`
	ContinueOnFailure bool      `json:"continue_on_failure"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Cron represents the cron timing for a schedule.
type Cron struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth string `json:"day_of_month"`
	Month      string `json:"month"`
	Hour       string `json:"hour"`
	Minute     string `json:"minute"`
}

// Schedule represents a schedule for running automated tasks on a server.
type Schedule struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Cron           Cron       `json:"cron"`
	IsActive       bool       `json:"is_active"`
	IsProcessing   bool       `json:"is_processing"`
	OnlyWhenOnline bool       `json:"only_when_online"`
	LastRunAt      *time.Time `json:"last_run_at"`
	NextRunAt      *time.Time `json:"next_run_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Tasks are populated by our SDK from the 'relationships' object on detail calls.
	Tasks []*Task `json:"-"`
}

// ScheduleCreateOptions defines the request body for creating a new schedule.
type ScheduleCreateOptions struct {
	Name           string `json:"name"`
	IsActive       *bool  `json:"is_active,omitempty"` // Use pointer to allow sending 'false'
	OnlyWhenOnline *bool  `json:"only_when_online,omitempty"`
	Minute         string `json:"minute"`
	Hour           string `json:"hour"`
	DayOfWeek      string `json:"day_of_week"`
	DayOfMonth     string `json:"day_of_month"`
	Month          string `json:"month"`
}

// ScheduleUpdateOptions defines the request body for updating an existing schedule.
// It is identical to the create options.
type ScheduleUpdateOptions = ScheduleCreateOptions

// TaskCreateOptions defines the request body for creating a new task.
type TaskCreateOptions struct {
	Action            string `json:"action"`
	Payload           string `json:"payload"`
	TimeOffset        int    `json:"time_offset"` // Must be >= 0
	ContinueOnFailure *bool  `json:"continue_on_failure,omitempty"`
}

// TaskUpdateOptions defines the request body for updating an existing task.
// It is identical to the create options.
type TaskUpdateOptions = TaskCreateOptions

// scheduleResponse and taskResponse are helpers for unmarshaling single object responses.
type ScheduleResponse struct {
	Object     string    `json:"object"`
	Attributes *Schedule `json:"attributes"`
}

type TaskResponse struct {
	Object     string `json:"object"`
	Attributes *Task  `json:"attributes"`
}

// scheduleDetailResponse is a special helper to unmarshal the schedule detail response,
// which includes the list of tasks in a 'relationships' object.
type ScheduleDetailResponse struct {
	Object        string    `json:"object"`
	Attributes    *Schedule `json:"attributes"`
	Relationships *struct {
		Tasks *PaginatedResponse[Task] `json:"tasks"`
	} `json:"relationships"`
}
