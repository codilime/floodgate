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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	cmd.SetArgs([]string{"--config=" + config, "apply", "-g"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Apply() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Apply() Read output err: %v", err)
	}

	t.Log(out)
	//
	//outStr := strings.TrimSpace(string(out))
	//if outStr != syncDryRunWant {
	//	t.Errorf("cmd.SyncDryRun() got:\n %s\n want:\n %s\n", outStr, syncDryRunWant)
	//}
}
