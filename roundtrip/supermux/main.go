package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/aunem/transpose/pkg/context"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var roundtripper http.RoundTripper

func main() {}

type superMux struct{}

// Spec represents the spec
var Spec SuperMuxSpec

// RoundtripPlugin is the roundtrip plugin inerface
var RoundtripPlugin superMux

func (s *superMux) Roundtrip(req context.Request) (context.Response, error) {
	var err error
	switch r := req.(type) {
	case *context.HTTPRequest:
		log.Debugf("executing http request: %+v", req.GetRequest())
		var match bool
		for _, v := range Spec.HTTP {
			u := r.Request.URL
			log.Debug("matching path...")
			match, err = path.Match(v.Path, u.Path)
			if err != nil {
				return nil, err
			}
			if match {
				log.Debug("path matched")
				// clean request
				rcc := RequestToRoundtrip(r)
				log.Debugf("rcc: %+v", rcc)
				host := fmt.Sprintf("%s:%s", v.Backend.ServiceName, v.Backend.ServicePort)
				u := r.Request.URL
				u.Host = host
				u.Scheme = "http"
				rcc.Request.Host = host
				changeTarget(rcc.Request, u)
				log.Debugf("change target: %+v", rcc.Request)
				resp, err := roundtripper.RoundTrip(rcc.Request)
				if err != nil {
					return nil, err
				}
				return &context.HTTPResponse{
					ID:       r.ID,
					Request:  r.Request,
					Response: resp,
				}, nil
			}
		}
		return nil, fmt.Errorf("path did not match")
	default:
		return nil, fmt.Errorf("request type uknown")
	}
}

func (s *superMux) LoadSpec(spec interface{}) error {
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

func (s *superMux) Init() error {
	roundtripper = http.DefaultTransport
	return nil
}

func (s *superMux) Stats() ([]byte, error) {
	return nil, nil
}
