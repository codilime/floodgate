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

func TestResource_SaveLocalState(t *testing.T) {
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
			got, err := r.SaveLocalState()
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

func TestResource_GetFullDiff(t *testing.T) {
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
		{
			name: "Empty remote json",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  emptyJSON,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: singleKeyJSON0Diff,
		},
		{
			name: "Empty local json",
			fields: fields{
				name:         "resource",
				localState:   emptyJSON,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: singleKeyJSON0Diff,
		},
		{
			name: "Proper diff twice single key",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  singleKeyJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: singleDiff,
		},
		{
			name: "Proper diff single key agains double key",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON1,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: oneKeyChanged01,
		},
		{
			name: "Proper diff twice double key",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON0,
				remoteState:  twoKeysJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: twoKeysChanged10,
		},
		{
			name: "Proper diff twice double key",
			fields: fields{
				name:         "resource",
				localState:   complexKeysJSON0110,
				remoteState:  complexKeysJSON1110,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: nestedChange,
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
			if got := r.GetFullDiff(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.GetFullDiff() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestResource_GetNormalizedDiff(t *testing.T) {
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
		{
			name: "Proper diff twice single key",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  singleKeyJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: singleDiff,
		},
		{
			name: "Proper diff single key agains double key",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON1,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: oneKeyChanged01,
		},
		{
			name: "Proper diff twice double key",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON0,
				remoteState:  twoKeysJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: twoKeysChanged10,
		},
		{
			name: "Proper diff single key complex key",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  complexKeysJSON1110,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: normalizedNewKey,
		},
		{
			name: "Proper diff twice double key",
			fields: fields{
				name:         "resource",
				localState:   complexKeysJSON0110,
				remoteState:  complexKeysJSON1110,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want: nestedChange,
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
			if got := r.GetNormalizedDiff(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.GetNormalizedDiff() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestResource_getNormalizedRemoteState(t *testing.T) {
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
			name: "Exact the same json twice",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    singleKeyJSON0,
			wantErr: false,
		},
		{
			name: "Same structure, different data",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  singleKeyJSON1,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    singleKeyJSON1,
			wantErr: false,
		},
		{
			name: "Local single key json",
			fields: fields{
				name:         "resource",
				localState:   singleKeyJSON0,
				remoteState:  twoKeysJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    singleKeyJSON0,
			wantErr: false,
		},
		{
			name: "Remote single key json",
			fields: fields{
				name:         "resource",
				localState:   twoKeysJSON0,
				remoteState:  singleKeyJSON0,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    singleKeyJSON0,
			wantErr: false,
		},
		{
			name: "Same complex structures, different data",
			fields: fields{
				name:         "resource",
				localState:   complexKeysJSON0110,
				remoteState:  complexKeysJSON1110,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    complexKeysJSON1110,
			wantErr: false,
		},
		{
			name: "Same complex structures, different data",
			fields: fields{
				name:         "resource",
				localState:   complexKeysJSON0110,
				remoteState:  complexKeysJSON1111,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    normalizedRemovedKeys,
			wantErr: false,
		},
		{
			name: "Broken local json",
			fields: fields{
				name:         "resource",
				localState:   brokenJSON,
				remoteState:  complexKeysJSON1111,
				spinnakerAPI: &gateclient.GateapiClient{},
			},
			want:    emptyJSON,
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
			got, err := r.getNormalizedRemoteState()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.getNormalizedRemoteState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.getNormalizedRemoteState() = %s, want %s", got, tt.want)
			}
		})
	}
}

var emptyJSON = []byte("{}")
var singleKeyJSON0 = []byte("{\"key1\":0}")
var singleKeyJSON1 = []byte("{\"key1\":1}")
var singleKeyJSON0Diff = []byte("{\n\"key1\": 0\n}")
var twoKeysJSON0 = []byte("{\"key1\":0, \"key2\": 0}")
var twoKeysJSON1 = []byte("{\"key1\":1, \"key2\": 1}")
var complexKeysJSON0110 = []byte("{\"key1\":{\"key2\":0,\"key3\":1},\"key4\":{\"key5\":1,\"key6\":0}}")
var complexKeysJSON1110 = []byte("{\"key1\":{\"key2\":1,\"key3\":1},\"key4\":{\"key5\":1,\"key6\":0}}")
var complexKeysJSON1111 = []byte("{\"key7\":{\"key8\":1,\"key9\":1},\"key4\":{\"key5\":1,\"key6\":1}}")
var brokenJSON = []byte("{")
var singleDiff = []byte("{\n\"key1\": 1 => 0\n}")
var oneKeyChanged01 = []byte("{\n\"key1\": 0 => 1,\n\"key2\": 1\n}")
var twoKeysChanged10 = []byte("{\n\"key1\": 1 => 0,\n\"key2\": 1 => 0\n}")
var newKeyAdded = []byte("{\n\"key1\": {} => 0,\n\"key4\": {\n\"key5\": 1,\n\"key6\": 0\n}\n}")
var nestedChange = []byte("{\n\"key1\": {\n\"key2\": 1 => 0,\n\"key3\": 1\n},\n\"key4\": {\n\"key5\": 1,\n\"key6\": 0\n}\n}")
var addingNestedKey = []byte("{\n\"key1\": {\n\"key2\": 0,\n\"key3\": 1\n},\n\"key4\": {\n\"key5\": 1,\n\"key6\": 1 => 0\n},\n\"key7\": {\n\"key8\": 1,\n\"key9\": 1\n}\n}")
var normalizedNewKey = []byte("{\n\"key1\": {} => 0\n}")
var normalizedKeyChanged = []byte("{\n\"key1\": {\n\"key2\": 0,\n\"key3\": 1\n},\n\"key4\": {\n\"key5\": 1,\n\"key6\": 1 => 0\n}\n}")
var normalizedRemovedKeys = []byte("{\"key4\":{\"key5\":1,\"key6\":1}}")
