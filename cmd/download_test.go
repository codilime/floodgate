package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/projects/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n  \"id\": \"b355c4ca-d70e-40d6-9505-8e21937b324f\",\n  \"name\": \"test\",\n  \"config\": {\n    \"pipelineConfigs\": [\n      {\n        \"application\": \"test\",\n        \"pipelineConfigId\": \"c6f10f23-be90-47af-854d-3a048c7875e3\"\n      }\n    ],\n    \"applications\": [\n      \"test\"\n    ]\n  }\n}"))
	})
	mux.HandleFunc("/applications/test/pipelineConfigs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[\n  {\n    \"keepWaitingPipelines\": false,\n    \"limitConcurrent\": true,\n    \"application\": \"test\",\n    \"lastModifiedBy\": \"test\",\n    \"name\": \"Deploy Application\",\n    \"stages\": [],\n    \"id\": \"c6f10f23-be90-47af-854d-3a048c7875e3\",\n    \"triggers\": []\n  }\n]"))
	})
	mux.HandleFunc("/applications/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n  \"name\": \"test\",\n  \"attributes\": {\n    \"name\": \"test\",\n    \"lastModifiedBy\": \"admin\",\n    \"cloudProviders\": \"kubernetes\",\n    \"trafficGuards\": [],\n    \"instancePort\": 80,\n    \"user\": \"admin\",\n    \"accounts\": \"spinnaker\"\n  },\n  \"clusters\": {}\n}"))
	})
	mux.HandleFunc("/applications/test/pipelineConfigs/Deploy Application", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\n    \"keepWaitingPipelines\": false,\n    \"limitConcurrent\": true,\n    \"application\": \"test\",\n    \"lastModifiedBy\": \"test\",\n    \"name\": \"Deploy Application\",\n    \"stages\": [],\n    \"id\": \"c6f10f23-be90-47af-854d-3a048c7875e3\",\n    \"triggers\": []\n  }"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.Download() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "download", "-p", "test", "-o", path.Join(dir, "downloadOutput")})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Download() Execute err: %v", err)
	}

	downloadDir := path.Join(dir, "downloadOutput", "test")
	applicationFiles, err := ioutil.ReadDir(path.Join(downloadDir, "applications"))
	if err != nil {
		t.Fatalf("cmd.Download() Read dir err: %v", err)
	}

	pipelineFiles, err := ioutil.ReadDir(path.Join(downloadDir, "pipelines"))
	if err != nil {
		t.Fatalf("cmd.Download() Read dir err: %v", err)
	}

	downloadApplicationFile := "test.json"
	downloadPipelineFile := "c6f10f23-be90-47af-854d-3a048c7875e3.json"

	for _, file := range applicationFiles {
		if file.Name() != downloadApplicationFile {
			t.Errorf("cmd.Download() Could not find %s in download/applications directory", downloadApplicationFile)
		}

		content, err := ioutil.ReadFile(path.Join(downloadDir, "applications", downloadApplicationFile))
		if err != nil {
			t.Fatalf("cmd.Download() Cannot read file %s err: %v", downloadApplicationFile, err)
		}

		contentStr := strings.TrimSpace(string(content))
		if contentStr != applicationWant {
			t.Errorf("cmd.Download() got:\n %s want:\n %s", contentStr, applicationWant)
		}
	}

	for _, file := range pipelineFiles {
		if file.Name() != downloadPipelineFile {
			t.Errorf("cmd.Download() Could not find %s in download/pipelines directory", downloadPipelineFile)
		}

		content, err := ioutil.ReadFile(path.Join(downloadDir, "pipelines", downloadPipelineFile))
		if err != nil {
			t.Fatalf("cmd.Download() Cannot read file %s err: %v", downloadPipelineFile, err)
		}

		contentStr := strings.TrimSpace(string(content))
		if contentStr != pipelineWant {
			t.Errorf("cmd.Download() got:\n %s want:\n %s", contentStr, pipelineWant)
		}
	}
}

var pipelineWant = "{\n\t\"application\": \"test\",\n\t\"id\": \"c6f10f23-be90-47af-854d-3a048c7875e3\",\n\t\"keepWaitingPipelines\": false,\n\t\"lastModifiedBy\": \"test\",\n\t\"limitConcurrent\": true,\n\t\"name\": \"Deploy Application\",\n\t\"stages\": [],\n\t\"triggers\": []\n}"
var applicationWant = "{\n\t\"accounts\": \"spinnaker\",\n\t\"cloudProviders\": \"kubernetes\",\n\t\"instancePort\": 80,\n\t\"lastModifiedBy\": \"admin\",\n\t\"name\": \"test\",\n\t\"trafficGuards\": [],\n\t\"user\": \"admin\"\n}"
