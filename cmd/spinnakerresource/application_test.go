package spinnakerresource

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/test"
)

func TestApplication_Init(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	type args struct {
		name      string
		localdata []byte
		ts        *httptest.Server
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with 200 OK response",
			args: args{
				name:      "app",
				localdata: []byte("{}"),
				ts:        test.MockGateServerReturn200(""),
			},
			wantErr: false,
		},
		{
			name: "with 500 ISE response",
			args: args{
				name:      "app",
				localdata: []byte("{}"),
				ts:        test.MockGateServerReturn500(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.args.ts
			defer ts.Close()
			api := test.MockGateapiClient(ts.URL)
			a := &Application{}
			err := a.Init(tt.args.name, api, tt.args.localdata)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.SaveLocalState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestApplication_loadRemoteState(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	testApp := fields{
		name: "app",
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
			remoteData: "{\"key\":1}",
			fields:     testApp,
			wantErr:    false,
		},
		{
			name:       "with 404 Not Found response",
			ts:         test.MockGateServerReturn404,
			remoteData: "{}",
			fields:     testApp,
			wantErr:    false,
		},
		{
			name:       "with 500 ISE response",
			ts:         test.MockGateServerReturn500,
			remoteData: "{}",
			fields:     testApp,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts(tt.remoteData)
			defer ts.Close()
			fmt.Println(ts.URL)
			api := test.MockGateapiClient(ts.URL)

			a := &Application{
				Resource: &Resource{
					spinnakerAPI: api,
				},
				name: tt.fields.name,
			}
			if err := a.loadRemoteState(); (err != nil) != tt.wantErr {
				t.Errorf("Application.loadRemoteState() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(a.GetRemoteState())
			if !tt.wantErr && (string(a.GetRemoteState()) != tt.remoteData) {
				t.Errorf("Data was loaded but not correctly stored: expected '%s', got '%s'", tt.remoteData, a.remoteState)
			}
		})
	}
}

func TestApplication_SaveRemoteState(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Application{
				Resource: &Resource{
					localState:   tt.fields.localState,
					remoteState:  tt.fields.remoteState,
					spinnakerAPI: tt.fields.spinnakerAPI,
				},
				name: tt.fields.name,
			}
			if err := a.SaveRemoteState(); (err != nil) != tt.wantErr {
				t.Errorf("Application.SaveRemoteState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
