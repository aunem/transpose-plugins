package main

import (
	"github.com/aunem/transpose/pkg/context"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type coralPlugin struct{}

// MiddlewarePlugin exports the plugin struct
var MiddlewarePlugin coralPlugin

// Spec exports the spec data
var Spec coralSpec

func main() {}

func (p *coralPlugin) ProcessRequest(req context.Request) (context.Request, error) {
	return nil, nil
}

func (p *coralPlugin) ProcessResponse(resp context.Response) (context.Response, error) {
	return nil, nil
}

func (p *coralPlugin) Init(spec interface{}) error {
	b, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &Spec)
	if err != nil {
		return err
	}
	log.Debugf("loaded spec: %+v", Spec)
	return nil
}

func (p *coralPlugin) Stats() ([]byte, error) {
	return nil, nil
}
