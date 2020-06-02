package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSyncDryRun(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL)
	if err != nil {
		t.Errorf("cmd.SyncDryRun() Error while creating temp config %v", err)
	}
	defer os.RemoveAll(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "sync", "-d"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.SyncDryRun() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.SyncDryRun() Read output err: %v", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr != syncDryRunWant {
		t.Errorf("cmd.SyncDryRun() got:\n %s\n want:\n %s\n", outStr, syncDryRunWant)
	}
}

var syncDryRunWant = "Following resources are changed:\nResource: jsonapp\nType: application\nChanges:\n@ []\n- null\n+ {\"cloudProviders\":\"kubernetes\",\"dataSources\":{\"disabled\":[],\"enabled\":[]},\"description\":\"Example application created from JSON file.\",\"email\":\"example@example.com\",\"name\":\"jsonapp\",\"platformHealthOnly\":false,\"platformHealthOnlyShowOverride\":false,\"providerSettings\":{\"aws\":{\"useAmiBlockDeviceMappings\":false},\"gce\":{\"associatePublicIpAddress\":false}},\"trafficGuards\":[],\"user\":\"floodgate@example.com\"}\n\nResource: jsonnetpipeline (Example pipeline from Jsonnet)\nType: pipeline\nChanges:\n@ []\n- null\n+ {\"application\":\"jsonnetapp\",\"id\":\"jsonnetpipeline\",\"keepWaitingPipelines\":false,\"limitConcurrent\":true,\"name\":\"Example pipeline from Jsonnet\",\"notifications\":[],\"stages\":[],\"triggers\":[]}\n\nResource: jsonnetpt (Example pipeline template from Jsonnet)\nType: pipelinetemplate\nChanges:\n@ []\n- null\n+ {\"id\":\"jsonnetpt\",\"metadata\":{\"description\":\"Example pipeline template created from Jsonnet file.\",\"name\":\"Example pipeline template from Jsonnet\",\"owner\":\"floodgate@example.com\",\"scopes\":[\"global\"]},\"protect\":false,\"schema\":\"v2\",\"variables\":[]}"
