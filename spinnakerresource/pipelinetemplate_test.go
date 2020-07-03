package spinnakerresource

import (
	"encoding/json"
	"testing"

	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/test"
)

func TestPipelineTemplate_Init(t *testing.T) {
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
			name: "server responds with 200 OK",
			args: args{
				localData: testPipelineTemplate,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: false,
		},
		{
			name: "server responds with 404 Not Found",
			args: args{
				localData: testPipelineTemplate,
				ts:        test.MockGateServerReturn404,
			},
			wantErr: false,
		},
		{
			name: "server responds with 500 ISE response",
			args: args{
				localData: testPipelineTemplate,
				ts:        test.MockGateServerReturn500,
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
			name: "missing key metadata",
			args: args{
				localData: PipelineTemplateWithMissingMetadata,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing key schema",
			args: args{
				localData: PipelineTemplateWithMissingSchema,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key name",
			args: args{
				localData: PipelineTemplateWithMissingMetadataName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key description",
			args: args{
				localData: PipelineTemplateWithMissingMetadataDescription,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key owner",
			args: args{
				localData: PipelineTemplateWithMissingMetadataOwner,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key scopes",
			args: args{
				localData: PipelineTemplateWithMissingMetadataScopes,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key id",
			args: args{
				localData: PipelineWithEmptyID,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key name",
			args: args{
				localData: PipelineTemplateWithEmptyMetadataName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key description",
			args: args{
				localData: PipelineTemplateWithEmptyMetadataDescription,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key owner",
			args: args{
				localData: PipelineTemplateWithEmptyMetadataOwner,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.args.ts("{}")
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			pt := &PipelineTemplate{}
			err := pt.Init(api, tt.args.localData)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestPipelineTemplate_LoadRemoteState(t *testing.T) {
	type fields struct {
		id           string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	pipelineTemplate := fields{
		id: "test-pipeline-template",
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
			fields:     pipelineTemplate,
			wantErr:    false,
		},
		{
			name:       "with 404 Not Found response",
			ts:         test.MockGateServerReturn404,
			remoteData: "{}",
			fields:     pipelineTemplate,
			wantErr:    false,
		},
		{
			name:       "with 500 ISE response",
			ts:         test.MockGateServerReturn500,
			remoteData: "{}",
			fields:     pipelineTemplate,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts(tt.remoteData)
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			pt := &PipelineTemplate{
				Resource: &Resource{},
				id:       tt.fields.id,
			}
			if err := pt.LoadRemoteState(api); (err != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", err, tt.wantErr)
			}
			if !tt.wantErr && (string(pt.GetRemoteState()) != tt.remoteData) {
				t.Errorf("data was loaded but not correctly stored: got %q, want %q", pt.remoteState, tt.remoteData)
			}
		})
	}
}

func TestPipelineTemplate_SaveRemoteState(t *testing.T) {
	type fields struct {
		id           string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	testPipelineTemplate := fields{
		id:         "test-pipeline-template",
		localState: []byte("{}"),
	}
	tests := []struct {
		name      string
		ts        test.MockGateServerFunction
		localData string
		fields    fields
		wantErr   bool
	}{
		{
			name:      "with 202 Accepted response",
			ts:        test.MockGateServerReturn202,
			localData: "{}",
			fields:    testPipelineTemplate,
			wantErr:   false,
		},
		{
			name:      "with 404 Not Found response",
			ts:        test.MockGateServerReturn404,
			localData: "{}",
			fields:    testPipelineTemplate,
			wantErr:   true,
		},
		{
			name:      "with 500 ISE response",
			ts:        test.MockGateServerReturn500,
			localData: "{}",
			fields:    testPipelineTemplate,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts("{}")
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			pt := PipelineTemplate{
				Resource: &Resource{
					localState:  tt.fields.localState,
					remoteState: tt.fields.remoteState,
				},
				id: tt.fields.id,
			}
			if err := pt.SaveLocalState(api); (err != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestPipelineTemplate_Name(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	pt := &PipelineTemplate{}
	err := pt.Init(api, testPipelineTemplate)
	if err != nil {
		t.Errorf("Resource.Name() error = %v", err)
	}

	want := "Test pipeline template"
	if pt.Name() != want {
		t.Errorf("Resource.Name() got %s, want %s", pt.Name(), want)
	}
}

func TestPipelineTemplate_ID(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	pt := &PipelineTemplate{}
	err := pt.Init(api, testPipelineTemplate)
	if err != nil {
		t.Errorf("Resource.ID() error = %v", err)
	}

	want := "test-pipeline-template"
	if pt.ID() != want {
		t.Errorf("Resource.ID() got %s, want %s", pt.ID(), want)
	}
}

func TestPipelineTemplate_LocalState(t *testing.T) {
	ts := test.MockGateServerReturn200("")
	api := test.MockGateapiClient(ts.URL)

	pt := &PipelineTemplate{}
	err := pt.Init(api, testPipelineTemplate)
	if err != nil {
		t.Errorf("Resource.LocalState() error = %v", err)
	}

	localState, _ := json.Marshal(testPipelineTemplate)
	if string(pt.LocalState()) != string(localState) {
		t.Errorf("Resource.LocalState() got %s, want %s", string(pt.LocalState()), string(localState))
	}
}

func TestPipelineTemplate_RemoteState(t *testing.T) {
	ts := test.MockGateServerReturn200("{}")
	api := test.MockGateapiClient(ts.URL)

	pt := &PipelineTemplate{}
	err := pt.Init(api, testPipelineTemplate)
	if err != nil {
		t.Errorf("Resource.RemoteState() error = %v", err)
	}

	if err := pt.LoadRemoteState(api); err != nil {
		t.Errorf("Resource.RemoteState() error = %v", err)
	}

	if string(pt.RemoteState()) != "{}" {
		t.Errorf("Resource.RemoteState() got %s, want %s", string(pt.RemoteState()), "{}")
	}
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

var PipelineTemplateWithMissingID = map[string]interface{}{
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

var PipelineTemplateWithMissingSchema = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
}

var PipelineTemplateWithMissingMetadata = map[string]interface{}{
	"id":     "test-pipeline-template",
	"schema": "v2",
}

var PipelineTemplateWithMissingMetadataName = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var PipelineTemplateWithMissingMetadataDescription = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":  "Test pipeline template",
		"owner": "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var PipelineTemplateWithMissingMetadataOwner = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var PipelineTemplateWithMissingMetadataScopes = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
	},
	"schema": "v2",
}

var PipelineTemplateWithEmptyID = map[string]interface{}{
	"id": "",
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

var PipelineTemplateWithEmptySchema = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "",
}

var PipelineTemplateWithEmptyMetadataName = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var PipelineTemplateWithEmptyMetadataDescription = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}

var PipelineTemplateWithEmptyMetadataOwner = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "",
		"scopes": []interface{}{
			"",
		},
	},
	"schema": "v2",
}
