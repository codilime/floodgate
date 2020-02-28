package spinnakerresource

import (
	"reflect"
	"testing"

	"github.com/codilime/floodgate/cmd/gateclient"
)

func TestResource_IsChanged(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Object is not changed",
			fields: fields{
				name:         "resource",
				localState:   emptyJSON,
				remoteState:  emptyJSON,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Object is changed",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON1,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Object is not changed (more keys)",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON1,
				remoteState:  twoKeysJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Object is changed (more keys)",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  twoKeysJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "localState is malformed",
			fields: fields{
				name:         "resource",
				localState:   brokenJSON,
				remoteState:  twoKeysJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "remoteState is malformed",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON1,
				remoteState:  brokenJSON,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Resource{
				name:         tt.fields.name,
				localState:   tt.fields.localState,
				remoteState:  tt.fields.remoteState,
				spinnakerAPI: tt.fields.spinnakerAPI,
			}
			got, err := r.IsChanged()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.IsChanged() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Resource.IsChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_GetLocalState(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "localState is returned",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  emptyJSON,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    singleKeyJSON0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Resource{
				name:         tt.fields.name,
				localState:   tt.fields.localState,
				remoteState:  tt.fields.remoteState,
				spinnakerAPI: tt.fields.spinnakerAPI,
			}
			got, err := r.GetLocalState()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.SaveLocalState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.SaveLocalState() = %v, want %v", got, tt.want)
			}
		})
	}
}

var emptyJSON = []byte("{}")
var singleKeyJSON0 = []byte("{\"key1\":0}")
var singleKeyJSON1 = []byte("{\"key1\":1}")
var twoKeysJSON1 = []byte("{\"key1\":1, \"key2\": 1}")
var brokenJSON = []byte("{")

func TestResource_GetRemoteState(t *testing.T) {
	type fields struct {
		name         string
		localState   []byte
		remoteState  []byte
		spinnakerAPI *gateclient.GateapiClient
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				name:         tt.fields.name,
				localState:   tt.fields.localState,
				remoteState:  tt.fields.remoteState,
				spinnakerAPI: tt.fields.spinnakerAPI,
			}
			if got := r.GetRemoteState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.GetRemoteState() = %v, want %v", got, tt.want)
			}
		})
	}
}
