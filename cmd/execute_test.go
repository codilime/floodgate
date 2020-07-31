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

func TestExecute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/webhooks/webhook/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"eventId\": \"123\",\"eventProcessed\": true}"))
	})
	mux.HandleFunc("/applications/jsonnetapp/pipelineConfigs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n  \"application\": \"jsonnetapp\",\n  \"id\": \"jsonnetpipeline\"\n}"))
	})
	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Execute() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "execute", "test"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Execute() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Execute() Read output err: %v", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr != executeWant {
		t.Errorf("cmd.Execute() got:\n %s want:\n %s", outStr, executeWant)
	}
}

func TestExecuteWait(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/webhooks/webhook/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"eventId\": \"123\",\"eventProcessed\": true}"))
	})
	mux.HandleFunc("/applications/*/executions/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[{\n  \"status\": \"SUCCEEDED\"}]"))
	})
	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Execute() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "execute", "test", "-w", "--wait-time=1"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Execute() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Execute() Read output err: %v", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr != executeWaitWant {
		t.Errorf("cmd.Execute() got:\n %s want:\n %s", outStr, executeWant)
	}
}

var executeWant = "triggering 'test'\nevent processed successfully\nexecution id is 123"
var executeWaitWant = "triggering 'test'\nevent processed successfully\nexecution id is 123\npipeline succeeded"
