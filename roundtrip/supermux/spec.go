package main

// SuperMuxSpec represnts a configuration for the super mux router
type SuperMuxSpec struct {
	HTTP []HTTPRouter `json:"http" yaml:"http"`
}

// HTTPRouter holds the route translation information
type HTTPRouter struct {
	Path    string      `json:"path" yaml:"path"`
	Backend HTTPBackend `json:"backend" yaml:"backend"`
}

// HTTPBackend represents the downstream service
type HTTPBackend struct {
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	ServicePort string `json:"servicePort" yaml:"servicePort"`
}
