package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aunem/transpose/config"
	"github.com/aunem/transpose/pkg/context"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

var hydraT *hydraPlugin
var token string

func TestMain(m *testing.M) {
	conf, err := config.LoadConfig("", "local")
	if err != nil {
		log.Fatal(err)
	}
	err = hydraT.LoadSpec(conf.Spec.Middleware.Request[0].Spec)
	if err != nil {
		log.Fatal(err)
	}
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("oryd/hydra", "latest", []string{"LOG_LEVEL=debug", "ISSUER=http://localhost:4444", "CONSENT_URL=http://localhost:3000/consent", "DATABASE_URL=memory", "FORCE_ROOT_CLIENT_CREDENTIALS=admin:password", "SYSTEM_SECRET=youReallyNeedToChangeThis"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		return hydraT.Init()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	tok, apiresp, err := hydraT.client.OAuth2Api.OauthToken()
	if err != nil {
		log.Fatal(err)
	}
	if apiresp.StatusCode >= 200 && apiresp.StatusCode <= 299 {
		log.Print("HTTP Status OK!")
	} else {
		log.Fatalf("status code non 2xx: %v from resp: %+v", apiresp.StatusCode, apiresp)
	}
	if tok.AccessToken == "" {
		log.Fatal("access token blank")
	}
	token = tok.AccessToken
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestProcessRequest(t *testing.T) {
	jsonStr := []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	basicPost, err := http.NewRequest(http.MethodPost, "http://test:8080/", bytes.NewBuffer(jsonStr))
	require.Nil(t, err)
	basicPost.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	advPost, err := http.NewRequest(http.MethodPost, "http://test:8080/adv", bytes.NewBuffer(jsonStr))
	require.Nil(t, err)
	advPost.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	missingPost, err := http.NewRequest(http.MethodPost, "http://test:8080/", bytes.NewBuffer(jsonStr))
	require.Nil(t, err)
	requesttests := []struct {
		in  context.Request
		out context.Request
		err bool
		msg string
	}{
		{
			in: &context.HTTPRequest{
				ID:      "myid",
				Request: basicPost,
				RW:      &httptest.ResponseRecorder{},
			},
			out: &context.HTTPRequest{
				ID:      "myid",
				Request: basicPost,
				RW:      &httptest.ResponseRecorder{},
			},
			err: false,
			msg: "should have allowed reqest",
		},
		{
			in: &context.HTTPRequest{
				ID:      "myid",
				Request: advPost,
				RW:      &httptest.ResponseRecorder{},
			},
			out: &context.HTTPRequest{
				ID:      "myid",
				Request: advPost,
				RW:      &httptest.ResponseRecorder{},
			},
			err: true,
			msg: "should not have allowed reqest, not map entry",
		},
		{
			in: &context.HTTPRequest{
				ID:      "myid",
				Request: missingPost,
				RW:      &httptest.ResponseRecorder{},
			},
			out: &context.HTTPRequest{
				ID:      "myid",
				Request: missingPost,
				RW:      &httptest.ResponseRecorder{},
			},
			err: true,
			msg: "should not have allowed reqest, no token",
		},
	}

	for _, tt := range requesttests {
		o, err := hydraT.ProcessRequest(tt.in)
		if tt.err {
			require.NotNil(t, err, tt.msg)
			continue
		} else {
			require.Nil(t, err, tt.msg)
		}
		require.Equal(t, tt.out, o, tt.msg)
	}
}

func TestProcessResponse(t *testing.T) {
}
