package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApply(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/tasks/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\": \"SUCCEEDED\"}"))
	})
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"ref\": \"test/test/test\"}"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.String())
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Apply() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "apply"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Apply() Execute err: %v", err)
	}

	_, err = ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Apply() Read output err: %v", err)
	}
}
