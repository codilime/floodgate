package spinnakerresource

import (
	"encoding/json"
	"reflect"

	"github.com/nsf/jsondiff"

	"github.com/codilime/floodgate/cmd/gateclient"
)

// Resource is basic struct for all spinnaker resources
type Resource struct {
	name                    string
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
func (r Resource) GetFullDiff() []byte {
	options := new(jsondiff.Options)
	_, diff := jsondiff.Compare(r.remoteState, r.localState, options)

	return []byte(diff)
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
func (r Resource) GetNormalizedDiff() []byte {
	options := new(jsondiff.Options)
	remoteState, _ := r.getNormalizedRemoteState()

	_, diff := jsondiff.Compare(remoteState, r.localState, options)

	return []byte(diff)
// GetRemoteState is used to view stored remote state.
func (r Resource) GetRemoteState() []byte {
	return r.remoteState
}

// Resourcer interface for Spinnaker resource
type Resourcer interface {
	// Init is used configure object and to load remote data into it
	Init() error
	// IsChanged is used to compare local and remmote state
	IsChanged() (bool, error)
	// SaveRemoteState is used to save state remotely
	SaveRemoteState() error
	// SaveLocalState is used to save state localy
	SaveLocalState() ([]byte, error)
}
