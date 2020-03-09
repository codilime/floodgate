package spinnakerresource

import (
	"testing"

	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/test"
)

func TestPipelineTemplate_Init(t *testing.T) {
	type args struct {
		id        string
		localData []byte
		ts        test.MockGateServerFunction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with 200 OK response",
			args: args{
				id:        "test-pipeline-template",
				localData: []byte("{}"),
				ts:        test.MockGateServerReturn200,
			},
			wantErr: false,
		},
		{
			name: "with 404 Not Found response",
			args: args{
				id:        "test-pipeline-template",
				localData: []byte("{}"),
				ts:        test.MockGateServerReturn404,
			},
			wantErr: false,
		},
		{
			name: "with 500 ISE response",
			args: args{
				id:        "test-pipeline-template",
				localData: []byte("{}"),
				ts:        test.MockGateServerReturn500,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.args.ts("")
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			pt := &PipelineTemplate{}
			err := pt.Init(tt.args.id, api, tt.args.localData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test case %#v: got error %v, wantErr: %v", tt.name, err, tt.wantErr)
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
				t.Errorf("Test case %#v: got error %v, wantErr: %v", tt.name, err, tt.wantErr)
			}
			if !tt.wantErr && (string(pt.GetRemoteState()) != tt.remoteData) {
				t.Errorf("Test case %#v: data was loaded but not correctly stored; got %q, want %q", tt.name, pt.remoteState, tt.remoteData)
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
			if err := pt.SaveRemoteState(); (err != nil) != tt.wantErr {
				t.Errorf("Test case %#v: got error %v, wantErr: %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
