package spinnakerresource

import (
	"encoding/json"
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
		localdata map[string]interface{}
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
				localdata: testAppLocalData,
				ts:        test.MockGateServerReturn200(""),
			},
			wantErr: false,
		},
		{
			name: "with 500 ISE response",
			args: args{
				name:      "app",
				localdata: testAppLocalData,
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
			err := a.Init(api, tt.args.localdata)
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
		remoteData map[string]interface{}
		fields     fields
		wantErr    bool
	}{
		{
			name:       "with 200 OK response",
			ts:         test.MockGateServerReturn200,
			remoteData: testAppRemoteData,
			fields:     testApp,
			wantErr:    false,
		},
		{
			name:       "with 404 Not Found response",
			ts:         test.MockGateServerReturn404,
			remoteData: map[string]interface{}{"attributes": map[string]interface{}{}},
			fields:     testApp,
			wantErr:    false,
		},
		{
			name:       "with 500 ISE response",
			ts:         test.MockGateServerReturn500,
			remoteData: testAppRemoteData,
			fields:     testApp,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsResponseData, _ := json.Marshal(tt.remoteData)
			ts := tt.ts(string(tsResponseData))
			defer ts.Close()
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
			remoteData := tt.remoteData["attributes"].(map[string]interface{})
			if _, exists := tt.remoteData["clusters"]; exists {
				clusters := tt.remoteData["clusters"].(map[string]interface{})
				if len(clusters) > 0 {
					remoteData["clusters"] = clusters
				}
			}
			remoteState, _ := json.Marshal(remoteData)
			if !tt.wantErr && (string(a.GetRemoteState()) != string(remoteState)) {
				t.Errorf("Application.loadRemoteState(): data was loaded but not correctly stored. Expected '%s', got '%s'", remoteState, a.remoteState)
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
			if err := a.SaveLocalState(); (err != nil) != tt.wantErr {
				t.Errorf("Application.SaveLocalState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testAppLocalData = map[string]interface{}{
	"name":           "testapplication",
	"description":    "Test application",
	"email":          "test@floodgate.com",
	"user":           "test@floodgate.com",
	"cloudProviders": "kubernetes",
	"dataSources": map[string]interface{}{
		"disabled": []string{},
		"enabled":  []string{},
	},
	"platformHealthOnly":             false,
	"platformHealthOnlyShowOverride": false,
	"providerSettings": map[string]interface{}{
		"aws": map[string]interface{}{
			"useAmiBlockDeviceMappings": false,
		},
		"gce": map[string]interface{}{
			"associatePublicIpAddress": false,
		},
	},
	"trafficGuards": []string{},
}

var testAppRemoteData = map[string]interface{}{
	"name": "test-application",
	"attributes": map[string]interface{}{
		"name":           "testapplication",
		"description":    "Test application",
		"email":          "test@floodgate.com",
		"user":           "test@floodgate.com",
		"cloudProviders": "kubernetes",
		"dataSources": map[string]interface{}{
			"disabled": []string{},
			"enabled":  []string{},
		},
		"platformHealthOnly":             false,
		"platformHealthOnlyShowOverride": false,
		"providerSettings": map[string]interface{}{
			"aws": map[string]interface{}{
				"useAmiBlockDeviceMappings": false,
			},
			"gce": map[string]interface{}{
				"associatePublicIpAddress": false,
			},
		},
		"trafficGuards": []string{},
	},
	"clusters": map[string]interface{}{},
}
