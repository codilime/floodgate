package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL)
	if err != nil {
		t.Errorf("cmd.Render() Error while creating temp config %v", err)
	}
	defer os.RemoveAll(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--verbose", "--config=" + config, "render", "-o", path.Join(dir, "renderOutput")})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Render() Execute err: %v", err)
	}

	renderDir := path.Join(dir, "renderOutput")
	files, err := ioutil.ReadDir(path.Join(renderDir, "pipelines"))
	if err != nil {
		t.Fatalf("cmd.Render() Read dir err: %v", err)
	}

	renderPipelineFile := "jsonnetpipeline.json"
	for _, file := range files {
		if file.Name() != renderPipelineFile {
			t.Errorf("cmd.Render() Could not find jsonnetpipeline.json in render directory")
		}

		content, err := ioutil.ReadFile(path.Join(renderDir, "pipelines", renderPipelineFile))
		if err != nil {
			t.Fatalf("cmd.Render() Cannot read file %s err: %v", renderPipelineFile, err)
		}

		contentStr := strings.TrimSpace(string(content))
		if contentStr != renderWant {
			t.Errorf("cmd.Render() got:\n %s want:\n %s", contentStr, inspectWant)
		}
	}
}

var renderWant = "{\n\t\"application\": \"jsonnetapp\",\n\t\"id\": \"jsonnetpipeline\",\n\t\"keepWaitingPipelines\": false,\n\t\"limitConcurrent\": true,\n\t\"name\": \"Example pipeline from Jsonnet\",\n\t\"notifications\": [],\n\t\"stages\": [],\n\t\"triggers\": []\n}"
