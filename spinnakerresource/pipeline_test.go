package spinnakerresource

import (
	"testing"

	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/test"
)

func TestPipeline_Init(t *testing.T) {
	type args struct {
		localData map[string]interface{}
		ts        test.MockGateServerFunction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "server responds with 200 OK and valid data",
			args: args{
				localData: testPipeline,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: false,
		},
		{
			name: "server responds with 200 OK and valid data with template",
			args: args{
				localData: testPipelineWithTemplate,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: false,
		},
		{
			name: "server responds with 404 Not Found",
			args: args{
				localData: testPipeline,
				ts:        test.MockGateServerReturn404,
			},
			wantErr: true,
		},
		{
			name: "server responds with 500 ISE response",
			args: args{
				localData: testPipeline,
				ts:        test.MockGateServerReturn500,
			},
			wantErr: true,
		},
		{
			name: "missing key name",
			args: args{
				localData: PipelineWithMissingName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing key application",
			args: args{
				localData: PipelineWithMissingApplication,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing key id",
			args: args{
				localData: PipelineWithMissingID,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty key name",
			args: args{
				localData: PipelineWithEmptyName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty key application",
			args: args{
				localData: PipelineWithEmptyApplication,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty key id",
			args: args{
				localData: PipelineWithEmptyID,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing schema for template",
			args: args{
				localData: PipelineWithMissingSchemaForTemplate,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty schema for template",
			args: args{
				localData: PipelineWithEmptySchemaForTemplate,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "template invalid type",
			args: args{
				localData: PipelineWithTemplateOfInvalidType,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.args.ts("{}")
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			p := &Pipeline{}
			err := p.Init(api, tt.args.localData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pipeline.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPipeline_LoadRemoteState(t *testing.T) {
	type fields struct {
		name         string
		application  string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	pipeline := fields{
		name:        "Test pipeline.",
		application: "Test pipeline application.",
	}
	tests := []struct {
		name       string
		ts         test.MockGateServerFunction
		remoteData string
		fields     fields
		wantErr    bool
	}{
		{
			name:       "with 200 OK response",
			ts:         test.MockGateServerReturn200,
			remoteData: "{}",
			fields:     pipeline,
			wantErr:    false,
		},
		{
			name:       "with 404 Not Found response",
			ts:         test.MockGateServerReturn404,
			remoteData: "{}",
			fields:     pipeline,
			wantErr:    true,
		},
		{
			name:       "with 500 ISE response",
			ts:         test.MockGateServerReturn500,
			remoteData: "{}",
			fields:     pipeline,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts(tt.remoteData)
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			p := &Pipeline{
				Resource:    &Resource{},
				name:        tt.fields.name,
				application: tt.fields.application,
			}
			if err := p.LoadRemoteState(api); (err != nil) != tt.wantErr {
				t.Errorf("Pipeline.loadRemoteState() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (string(p.GetRemoteState()) != tt.remoteData) {
				t.Errorf("data was loaded but not correctly stored: got %q, want %q", p.remoteState, tt.remoteData)
			}
		})
	}
}

func TestPipeline_SaveLocalState(t *testing.T) {
	type fields struct {
		name         string
		application  string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	testPipeline := fields{
		name:        "Test pipeline.",
		application: "Test pipeline application.",
		localState:  []byte("{}"),
	}
	tests := []struct {
		name      string
		ts        test.MockGateServerFunction
		localData string
		fields    fields
		wantErr   bool
	}{
		{
			name:      "with 200 OK response",
			ts:        test.MockGateServerReturn200,
			localData: "{}",
			fields:    testPipeline,
			wantErr:   false,
		},
		{
			name:      "with 404 Not Found response",
			ts:        test.MockGateServerReturn404,
			localData: "{}",
			fields:    testPipeline,
			wantErr:   true,
		},
		{
			name:      "with 500 ISE response",
			ts:        test.MockGateServerReturn500,
			localData: "{}",
			fields:    testPipeline,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts("{}")
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			p := Pipeline{
				Resource: &Resource{
					localState:  tt.fields.localState,
					remoteState: tt.fields.remoteState,
				},
				name:        tt.fields.name,
				application: tt.fields.application,
			}
			if err := p.SaveLocalState(api); (err != nil) != tt.wantErr {
				t.Errorf("Pipeline.SaveLocalState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testPipeline = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
}

var testPipelineWithTemplate = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
	"metadata": map[string]interface{}{
		"schema": "v2",
	},
}

var PipelineWithEmptyName = map[string]interface{}{
	"name":        "",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
}

var PipelineWithEmptyApplication = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "",
	"id":          "testpipeline",
}

var PipelineWithEmptyID = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "",
}

var PipelineWithMissingName = map[string]interface{}{
	"application": "Test pipeline application.",
	"id":          "testpipeline",
}

var PipelineWithMissingApplication = map[string]interface{}{
	"name": "Test pipeline.",
	"id":   "testpipeline",
}

var PipelineWithMissingID = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
}

var PipelineWithMissingSchemaForTemplate = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
	"template":    map[string]interface{}{"": ""},
}

var PipelineWithEmptySchemaForTemplate = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
	"template": map[string]interface{}{
		"schema": "",
	},
}

var PipelineWithTemplateOfInvalidType = map[string]interface{}{
	"name":        "Test pipeline.",
	"application": "Test pipeline application.",
	"id":          "testpipeline",
	"template":    "definitely not a map[string]interface{}",
}
