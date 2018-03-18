package main

import (
	"github.com/aunem/transpose/pkg/context"
)

type authPlugin struct{}

// Plugin exports

// Spec exports the config spec
var Spec HydraSpec

// MiddlewarePlugin exports the plugin interface
var MiddlewarePlugin authPlugin

func main() {}

func (m *authPlugin) ProcessRequest(req context.Request, spec interface{}) (context.Request, error) {
	return nil, nil
}

func (m *authPlugin) ProcessResponse(resp context.Response, spec interface{}) (context.Response, error) {
	return nil, nil
}

func (m *authPlugin) Init(spec interface{}) error {
	return nil
}

func (m *authPlugin) Stats() ([]byte, error) {
	return nil, nil
}
