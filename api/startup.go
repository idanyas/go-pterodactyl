package api

type StartupVariable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	EnvVariable  string `json:"env_variable"`
	DefaultValue string `json:"default_value"`
	ServerValue  string `json:"server_value"`
	IsEditable   bool   `json:"is_editable"`
	Rules        string `json:"rules"`
}

type UpdateVariableOptions struct {
	Key string `json:"key"`

	Value string `json:"value"`
}

type UpdateVariableResponse struct {
	Object     string           `json:"object"`
	Attributes *StartupVariable `json:"attributes"`
}
