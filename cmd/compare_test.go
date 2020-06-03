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

func TestCompare(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Compare() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "compare"})

	err = cmd.Execute()
	if err != nil && err.Error() != "end diff" {
		t.Fatalf("cmd.Compare() Execute err: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Compare() Read output err: %v", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr != compareWant {
		t.Errorf("cmd.Compare() got: %s want: %s", outStr, compareWant)
	}
}

var compareWant = "jsonapp (application)\n@ []\n- null\n+ {\"cloudProviders\":\"kubernetes\",\"dataSources\":{\"disabled\":[],\"enabled\":[]},\"description\":\"Example application created from JSON file.\",\"email\":\"example@example.com\",\"name\":\"jsonapp\",\"platformHealthOnly\":false,\"platformHealthOnlyShowOverride\":false,\"providerSettings\":{\"aws\":{\"useAmiBlockDeviceMappings\":false},\"gce\":{\"associatePublicIpAddress\":false}},\"trafficGuards\":[],\"user\":\"floodgate@example.com\"}\n\njsonnetpipeline (Example pipeline from Jsonnet) (pipeline)\n@ []\n- null\n+ {\"application\":\"jsonnetapp\",\"id\":\"jsonnetpipeline\",\"keepWaitingPipelines\":false,\"limitConcurrent\":true,\"name\":\"Example pipeline from Jsonnet\",\"notifications\":[],\"stages\":[],\"triggers\":[]}\n\njsonnetpt (Example pipeline template from Jsonnet) (pipelinetemplate)\n@ []\n- null\n+ {\"id\":\"jsonnetpt\",\"metadata\":{\"description\":\"Example pipeline template created from Jsonnet file.\",\"name\":\"Example pipeline template from Jsonnet\",\"owner\":\"floodgate@example.com\",\"scopes\":[\"global\"]},\"protect\":false,\"schema\":\"v2\",\"variables\":[]}"
