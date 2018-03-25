package main

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/aunem/transpose/pkg/context"
	"github.com/aunem/transpose/pkg/middleware"
	"github.com/aunem/transpose/pkg/roundtrip"
	log "github.com/sirupsen/logrus"
)

// Spec represents the plugins spec
var Spec HTTPListenerSpec

// ListenerPlugin is an export for the listener interface
var ListenerPlugin httpListener

type httpListener struct{}

func main() {}

// Listen implements the listener plugin inerface
func (h *httpListener) Listen(mw *middleware.Manager, rt *roundtrip.Manager) error {

	t := HTTPTranslator{
		FlushInterval: 10 * time.Second,
	}
	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// process request
		log.Debugf("processing request: %+v", req)
		rc := context.NewHTTPRequest(req)
		rcf, err := mw.ExecRequestStack(rc)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		// send to backend
		log.Debugf("sending roundtrip: %+v", rcf)
		rs, err := rt.ExecRoundtrip(rcf)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		// process response
		log.Debugf("processing response: %+v", rs)
		rsf, err := mw.ExecResponseStack(rs)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		log.Debugf("writing response: %+v", rsf)
		r, err := context.ResponseToHTTP(rsf)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}
		rw = t.ResponseToWriter(r.Response, rw)
	})
	port := Spec.Port
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
	log.Infof("starting server: %+v", s)
	return s.ListenAndServe()
}

// LoadSpec implements the listener plugin interface for loading the spec config
func (h *httpListener) LoadSpec(spec interface{}) error {
	log.Debug("loading spec...")
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

// Stats implements the listener plugin interface for fetching stats
func (h *httpListener) Stats() ([]byte, error) { return nil, nil }
