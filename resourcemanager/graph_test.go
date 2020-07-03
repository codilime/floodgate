package resourcemanager

import (
	"github.com/codilime/floodgate/gateclient"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/codilime/floodgate/test"
	"strings"
	"testing"
)

func initResourceGraph(api *gateclient.GateapiClient) (*ResourceGraph, error) {
	a := &spr.Application{}
	err := a.Init(api, testApplication)
	if err != nil {
		return nil, err
	}

	p := &spr.Pipeline{}
	err = p.Init(api, testPipeline)
	if err != nil {
		return nil, err
	}

	pt := &spr.PipelineTemplate{}
	err = pt.Init(api, testPipelineTemplate)
	if err != nil {
		return nil, err
	}

	resourceGraph := ResourceGraph{
		Resources: spr.SpinnakerResources{
			Applications: []*spr.Application{
				a,
			},
			Pipelines: []*spr.Pipeline{
				p,
			},
			PipelineTemplates: []*spr.PipelineTemplate{
				pt,
			},
		},
	}

	return &resourceGraph, nil
}

func TestResourceGraph_Create(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	resourceGraph, err := initResourceGraph(api)
	if err != nil {
		t.Errorf("ResourceGraph.Create() got error %v", err)
	}
	resourceGraph.Create()

	graphStr := strings.TrimSpace(resourceGraph.Graph.String())
	if graphStr != createWant {
		t.Errorf("ResourceGraph.Create() got %s, want %s", graphStr, createWant)
	}
}

func TestResourceGraph_Apply(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	resourceGraph, err := initResourceGraph(api)
	if err != nil {
		t.Errorf("ResourceGraph.Apply() got error %v", err)
	}
	resourceGraph.Create()

	err = resourceGraph.Apply(api, 5)
	if err != nil {
		t.Errorf("ResourceGraph.Apply() got error %v", err)
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
	"schema":      "v2",
	"template": map[string]interface{}{
		"reference": "spinnaker://test-pipeline-template",
	},
}

var testPipelineTemplate = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var createWant = "Start\nTest pipeline template\n  Start\nTest pipeline.\n  Test pipeline template\n  testapplication\ntestapplication\n  Start"
