package main

import (
	"fmt"
	"path"

	"github.com/aunem/transpose/pkg/context"
	"github.com/ory/fosite"
	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type hydraPlugin struct {
	client *hydra.CodeGenSDK
}

// Plugin exports

// Spec exports the config spec
var Spec HydraSpec

// MiddlewarePlugin exports the plugin interface
var MiddlewarePlugin hydraPlugin

func main() {}

func (m *hydraPlugin) ProcessRequest(req context.Request) (context.Request, error) {
	switch r := req.(type) {
	case *context.HTTPRequest:
		for _, v := range Spec.ResourceMap {
			match, err := path.Match(v.HTTP.Path, r.Request.URL.Path)
			if err != nil {
				return nil, err
			}
			if v.HTTP.Method != r.Request.Method {
				match = false
			}
			if match {
				log.Debugf("map match: %+v", v)
				token := fosite.AccessTokenFromRequest(r.Request)
				wr := swagger.WardenTokenAccessRequest{
					Resource: v.Resource,
					Action:   v.Action,
					Token:    token,
				}
				log.Debug("warden request: %+v", wr)
				resp, apiresp, err := m.client.DoesWardenAllowTokenAccessRequest(wr)
				if err != nil {
					return nil, err
				}
				if apiresp.StatusCode >= 200 && apiresp.StatusCode <= 299 {
					log.Debug("HTTP Status OK!")
				} else {
					return nil, fmt.Errorf("status code non 2xx: %v from resp: %+v", apiresp.StatusCode, apiresp)
				}
				if !resp.Allowed {
					return nil, fmt.Errorf("request not allowed")
				}
				log.Debug("request allowed")
				return req, nil
			}

		}
		return nil, fmt.Errorf("no resource path match")
	default:
		return nil, fmt.Errorf("request type unknown: %+v", r)
	}
}

func (m *hydraPlugin) ProcessResponse(resp context.Response) (context.Response, error) {
	return nil, nil
}

func (m *hydraPlugin) LoadSpec(spec interface{}) error {
	b, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &Spec)
	if err != nil {
		return err
	}
	log.Debugf("loaded spec: %+v", Spec)
	m.Init()
	return nil
}

func (m *hydraPlugin) Init() error {
	h := &hydra.Configuration{
		EndpointURL:  Spec.Endpoint,
		ClientID:     Spec.ClientID,
		ClientSecret: Spec.ClientSecret, // TODO: need better secrets injection
		Scopes:       Spec.Scopes,
	}
	s, err := hydra.NewSDK(h)
	if err != nil {
		return err
	}
	m.client = s
	return nil
}

func (m *hydraPlugin) Stats() ([]byte, error) {
	return nil, nil
}
