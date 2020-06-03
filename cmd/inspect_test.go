package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInspect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/applications/jsonapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"name\": \"jsonapp\",\"attributes\": {\"name\": \"jsonapp\"},\"clusters\": {}}"))
	})
	mux.HandleFunc("/applications/jsonnetapp/pipelineConfigs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n  \"application\": \"jsonnetapp\",\n  \"id\": \"jsonnetpipeline\"\n}"))
	})
	mux.HandleFunc("/v2/pipelineTemplates/jsonnetpt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n  \"application\": \"jsonnetpt\",\n  \"id\": \"jsonnetpt\"\n}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Inspect() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "inspect"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Inspect() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Inspect() Read output err: %v", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr != inspectWant {
		t.Errorf("cmd.Inspect() got:\n %s want:\n %s", outStr, inspectWant)
	}
}

var inspectWant = "Current Spinnaker resource status:\n\nApplications:\n{\"name\":\"jsonapp\"}\n\nPipelines:\n{\"application\":\"jsonnetapp\",\"id\":\"jsonnetpipeline\"}\n\nPipeline templates:\n{\"application\":\"jsonnetpt\",\"id\":\"jsonnetpt\"}"
