package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aunem/transpose/config"
	"github.com/aunem/transpose/pkg/context"
	"github.com/aunem/transpose/pkg/middleware"
	"github.com/aunem/transpose/pkg/roundtrip"
	log "github.com/sirupsen/logrus"
)

// Plugin exports
var Spec HTTPListenerSpec
var ListenerPlugin HTTPListenerSpec

// test
type HTTPListener struct{}

func main() {}

// Listen implements the listener plugin inerface
func (h *HTTPListener) Listen(spec config.TransposeSpec) error {
	httpSpec, ok := spec.Listener.Spec.(HTTPListenerSpec)
	if !ok {
		return fmt.Errorf("could not cast spec")
	}

	log.Info("printing...")
	log.Debugf("creating middleware manager....")
	mw, err := middleware.NewManager(spec)
	if err != nil {
		log.Fatalf("could not create middleware: %v", err)
	}
	log.Debugf("creating roundtrip manager...")
	rt, err := roundtrip.NewManager(spec)
	if err != nil {
		log.Fatalf("could not create roundtrip: %v", err)
	}
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
	if httpSpec.Port == "" {
		httpSpec.Port = "8080"
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", spec.Port),
		Handler: handler,
	}
	return s.ListenAndServe()
}

// Stats implements the listener plugin interface for fetching stats
func (h *HTTPListener) Stats() ([]byte, error) { return nil, nil }
