package test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/codilime/floodgate/cmd/gateclient"
	gateapi "github.com/codilime/floodgate/gateapi"
)

// MockGateServerFunction is a handler to a function for Gate server mock
type MockGateServerFunction func(string) *httptest.Server

// MockGateapiClient creates a basic API client without authentication.
func MockGateapiClient(gateURL string) *gateclient.GateapiClient {
	cfg := gateapi.NewConfiguration()
	cfg.BasePath = gateURL
	client := gateapi.NewAPIClient(cfg)

	return &gateclient.GateapiClient{
		APIClient: client,
		Context:   context.Background(),
	}
}

// MockGateServerReturn200 creates a HTTP server which returns code 200 and data.
func MockGateServerReturn200(data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	}))
}

// MockGateServerReturn404 creates a HTTP server which returns code 404 and data.
func MockGateServerReturn404(data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	}))
}

// MockGateServerReturn500 creates a HTTP server which returns code 500 and data.
func MockGateServerReturn500(data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
}
