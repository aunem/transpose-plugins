package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/aunem/transpose/pkg/context"
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
		var match bool
		for _, v := range Spec.HTTP {
			u := r.Request.URL
			match, err = path.Match(v.Path, u.Path)
			if err != nil {
				return nil, err
			}
			if match {
				// clean request
				rcc := RequestToRoundtrip(r)
				u := r.Request.URL
				host := fmt.Sprintf("%s:%s", v.Backend.ServiceName, v.Backend.ServicePort)
				u.Host = host
				changeTarget(rcc.Request, u)
				resp, err := roundtripper.RoundTrip(r.Request)
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
	return nil
}

func (s *superMux) Init() error {
	roundtripper = http.DefaultTransport
	return nil
}

func (s *superMux) Stats() ([]byte, error) {
	return nil, nil
}
