package spinnakerresource

import (
	"encoding/json"
	"reflect"

	jd "github.com/josephburnett/jd/lib"

	"github.com/codilime/floodgate/gateclient"
)

// Resource is basic struct for all spinnaker resources
type Resource struct {
	localState, remoteState []byte
	spinnakerAPI            *gateclient.GateapiClient
}

// IsChanged is used to compare local and remmote state
func (r Resource) IsChanged() (bool, error) {
	var (
		localJSON, remoteJSON map[string]interface{}
	)

	if err := json.Unmarshal(r.localState, &localJSON); err != nil {
		return false, err
	}
	if err := json.Unmarshal(r.remoteState, &remoteJSON); err != nil {
		return false, err
	}
	for k := range localJSON {
		if _, exists := remoteJSON[k]; !exists {
			return true, nil
		}
		if !reflect.DeepEqual(localJSON[k], remoteJSON[k]) {
			return true, nil
		}
	}

	return false, nil
}

// GetLocalState returns local state of an object.
func (r Resource) GetLocalState() ([]byte, error) {
	return r.localState, nil
}

// GetFullDiff function returns diff against remote state
func (r Resource) GetFullDiff() string {
	a, _ := jd.ReadJsonString(string(r.remoteState))
	b, _ := jd.ReadJsonString(string(r.localState))

	return a.Diff(b).Render()
}

func (r Resource) getNormalizedRemoteState() ([]byte, error) {
	var (
		localJSON, remoteJSON map[string]interface{}
	)
	if err := json.Unmarshal(r.localState, &localJSON); err != nil {
		return []byte("{}"), err
	}

	if err := json.Unmarshal(r.remoteState, &remoteJSON); err != nil {
		return []byte("{}"), err
	}

	remoteJSONNormalized := make(map[string]interface{})
	for k := range localJSON {
		if key, exists := remoteJSON[k]; exists {
			remoteJSONNormalized[k] = key
		}
	}

	remoteJSONNormalizedByte, err := json.Marshal(remoteJSONNormalized)
	if err != nil {
		return []byte("{}"), err
	}

	return remoteJSONNormalizedByte, nil
}

// GetNormalizedDiff function returns diff on only managed resources
func (r Resource) GetNormalizedDiff() string {
	remoteState, _ := r.getNormalizedRemoteState()
	a, _ := jd.ReadJsonString(string(remoteState))
	b, _ := jd.ReadJsonString(string(r.localState))

	return a.Diff(b).Render()
}

// GetRemoteState is used to view stored remote state.
func (r Resource) GetRemoteState() []byte {
	return r.remoteState
}

// Resourcer interface for Spinnaker resource
type Resourcer interface {
	// Init is used configure object and to load remote data into it
	Init(api *gateclient.GateapiClient, localData map[string]interface{}) error
	// IsChanged is used to compare local and remmote state
	IsChanged() (bool, error)
	// LocalState get resource's local state
	LocalState() []byte
	// RemoteState get resource's remote state
	RemoteState() []byte
	// SaveLocalState save resource's local state to Spinnaker
	SaveLocalState() error
}
