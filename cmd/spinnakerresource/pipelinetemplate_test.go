package spinnakerresource

import (
	"testing"

	"github.com/codilime/floodgate/cmd/gateclient"
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
				localData: missingID,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing key metadata",
			args: args{
				localData: missingMetadata,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing key schema",
			args: args{
				localData: missingSchema,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key name",
			args: args{
				localData: missingMetadataName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key description",
			args: args{
				localData: missingMetadataDescription,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key owner",
			args: args{
				localData: missingMetadataOwner,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "missing metadata key scopes",
			args: args{
				localData: missingMetadataScopes,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key id",
			args: args{
				localData: emptyID,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key name",
			args: args{
				localData: emptyMetadataName,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key description",
			args: args{
				localData: emptyMetadataDescription,
				ts:        test.MockGateServerReturn200,
			},
			wantErr: true,
		},
		{
			name: "empty metadata key owner",
			args: args{
				localData: emptyMetadataOwner,
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

func TestPipelineTemplate_loadRemoteState(t *testing.T) {
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
				Resource: &Resource{
					spinnakerAPI: api,
				},
				id: tt.fields.id,
			}
			if err := pt.loadRemoteState(); (err != nil) != tt.wantErr {
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
					localState:   tt.fields.localState,
					remoteState:  tt.fields.remoteState,
					spinnakerAPI: api,
				},
				id: tt.fields.id,
			}
			if err := pt.SaveLocalState(); (err != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", err, tt.wantErr)
			}
		})
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

var missingID = map[string]interface{}{
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

var missingSchema = map[string]interface{}{
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

var missingMetadata = map[string]interface{}{
	"id":     "test-pipeline-template",
	"schema": "v2",
}

var missingMetadataName = map[string]interface{}{
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

var missingMetadataDescription = map[string]interface{}{
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

var missingMetadataOwner = map[string]interface{}{
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

var missingMetadataScopes = map[string]interface{}{
	"id": "test-pipeline-template",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
	},
	"schema": "v2",
}

var emptyID = map[string]interface{}{
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

var emptySchema = map[string]interface{}{
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

var emptyMetadataName = map[string]interface{}{
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

var emptyMetadataDescription = map[string]interface{}{
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

var emptyMetadataOwner = map[string]interface{}{
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
