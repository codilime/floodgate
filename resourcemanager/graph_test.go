package resourcemanager

import (
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/codilime/floodgate/test"
	"strings"
	"testing"
)

func TestResourceGraph_Create(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	a := &spr.Application{}
	err := a.Init(api, testApplication)
	if err != nil {
		t.Errorf("ResourceGraph.Create() error %v\"", err)
	}

	p := &spr.Pipeline{}
	err = p.Init(api, testPipeline)
	if err != nil {
		t.Errorf("ResourceGraph.Create() error %v", err)
	}

	resourceGraph := ResourceGraph{
		Resources: spr.SpinnakerResources{
			Applications: []*spr.Application{
				a,
			},
			Pipelines: []*spr.Pipeline{
				p,
			},
		},
	}
	resourceGraph.Create()

	graphStr := strings.TrimSpace(resourceGraph.Graph.String())
	if graphStr != createWant {
		t.Errorf("ResourceGraph.Create() got %s, want %s", graphStr, createWant)
	}
}

var testApplication = map[string]interface{}{
	"name":  "testapplication",
	"email": "test@floodgate.com",
}

var testPipeline = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "testapplication",
	"id":          "testpipeline",
}

var createWant = "Start\nTest pipeline.\n  testapplication\ntestapplication\n  Start"
