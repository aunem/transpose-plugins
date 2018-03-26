package main

// HydraSpec represents the config spec
type HydraSpec struct {
	ClientID     string        `json:"clientId" yaml:"clientId"`
	ClientSecret string        `json:"clientSecret" yaml:"clientSecret"`
	Endpoint     string        `json:"endpoint" yaml:"endpoint"`
	Scopes       []string      `json:"scopes" yaml:"scopes"`
	ResourceMap  []Correlation `json:"resourceMap" yaml:"resourceMap"`
}

// Correlation correlates a request to a resource type
type Correlation struct {
	HTTP     HTTPReq `json:"http" yaml:"http"`
	Action   string  `json:"action" yaml:"action"`
	Resource string  `json:"resource" yaml:"resource"`
}

// HTTPReq represents a simple http request
type HTTPReq struct {
	Path   string `json:"path" yaml:"path"`
	Method string `json:"method" yaml:"method"`
}
