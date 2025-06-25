package api

import "time"

type Egg struct {
	ID           int               `json:"id"`
	UUID         string            `json:"uuid"`
	NestID       int               `json:"nest_id"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	DockerImages map[string]string `json:"docker_images"`
	Config       any               `json:"config"`
	Startup      any               `json:"startup"`
	Script       any               `json:"script"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type EggConfig struct {
	Files   any    `json:"files"`
	Startup any    `json:"startup"`
	Stop    string `json:"stop"`
	Logs    any    `json:"logs"`
}

type EggStartup struct {
	Done            string   `json:"done"`
	UserInteraction []string `json:"user_interaction"`
}

type EggScript struct {
	Privileged bool   `json:"privileged"`
	Install    string `json:"install"`
	Entry      string `json:"entry"`
	Container  string `json:"container"`
	Extends    any    `json:"extends"`
}
