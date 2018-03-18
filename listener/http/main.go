package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aunem/transpose/pkg/context"
	"github.com/aunem/transpose/pkg/middleware"
	"github.com/aunem/transpose/pkg/roundtrip"
	log "github.com/sirupsen/logrus"
)

// Plugin exports
var Spec HTTPListenerSpec
var ListenerPlugin HTTPListenerSpec

type HTTPListener {}

func main() {}

// Listen implements the listener plugin inerface
func (h *HTTPListener) Listen(spec interface{}) error {
	spec, ok := spec.(HTTPListenerSpec)
	if !ok {
		return fmt.Errorf("could not cast spec")
	}

	log.Debugf("creating middleware manager....")
	mw, err := middleware.NewManager(conf)
	if err != nil {
		log.Fatalf("could not create middleware: %v", err)
	}
	log.Debugf("creating roundtrip manager...")
	rt, err := roundtrip.NewManager(conf)
	if err != nil {
		log.Fatalf("could not create roundtrip: %v", err)
	}
	t := HTTPTranslator{
		FlushInterval: 10 * time.Second,
	}
	h := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// process request
		log.Debugf("processing request: %+v", req)
		rc := context.NewRequestContext(req)
		rcf, err := mw.ExecRequestStack(rc)
		if err != nil {
			log.Error(err)
			// need to bail
		}

		// send to backend
		log.Debugf("sending roundtrip: %+v", rcc)
		rs, err := rt.ExecRoundtrip(rcc)
		if err != nil {
			log.Error(err)
		}

		// process response
		log.Debugf("processing response: %+v", rs)
		rsf, err := mw.ExecResponseStack(rs)
		if err != nil {
			log.Error(err)
		}

		// write response
		log.Debugf("writing response: %+v", rsf)
		rw = t.ResponseToWriter(rsf.Response, rw)
	})
	if conf.Port == "" {
		conf.Port = "8080"
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.Port),
		Handler: h,
	}
	s.ListenAndServe()
}

func (h *HTTPListener) Stats() ([]byte, error) {}