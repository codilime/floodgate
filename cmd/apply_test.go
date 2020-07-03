package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
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

func TestApplyGraph(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/loggedOut", test.MockGateServerAuthLoggedOutHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	ts := httptest.NewServer(mux)

	dir, config, err := CreateTempFiles(ts.URL, false)
	if err != nil {
		t.Errorf("cmd.ApplyGraph() Error while creating temp config %v", err)
	}
	defer RemoveTempDir(dir)

	graphFileName := "graph.dot"

	b := bytes.NewBufferString("")
	cmd := NewRootCmd(b)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--config=" + config, "apply", "-dot", "-o", path.Join(dir, graphFileName)})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.ApplyGraph() Execute err: %v", err)
	}

	_, err = ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.ApplyGraph() Read output err: %v", err)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("cmd.ApplyGraph() Read dir err: %v", err)
	}

	for _, file := range files {
		if file.Name() == graphFileName {
			content, err := ioutil.ReadFile(path.Join(dir, graphFileName))
			if err != nil {
				t.Fatalf("cmd.ApplyGraph() Cannot read file %s err: %v", graphFileName, err)
			}

			if !bytes.Equal(content, graphFileWant) {
				t.Errorf("cmd.ApplyGraph() dot file content is incorrect")
			}
		}
	}
}

var graphFileWant = []byte{100, 105, 103, 114, 97, 112, 104, 32, 123, 10, 9, 99, 111, 109, 112, 111, 117, 110, 100, 32, 61, 32, 34, 116, 114, 117, 101, 34, 10, 9, 110, 101, 119, 114, 97, 110, 107, 32, 61, 32, 34, 116, 114, 117, 101, 34, 10, 9, 115, 117, 98, 103, 114, 97, 112, 104, 32, 34, 114, 111, 111, 116, 34, 32, 123, 10, 9, 9, 34, 91, 114, 111, 111, 116, 93, 32, 106, 115, 111, 110, 97, 112, 112, 34, 32, 45, 62, 32, 34, 91, 114, 111, 111, 116, 93, 32, 83, 116, 97, 114, 116, 34, 10, 9, 125, 10, 125, 10}
