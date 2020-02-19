package spinnakerresource

import (
	"encoding/json"
	"reflect"

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
	normalizedLocalState, err := r.normalizeJSON(r.localState)
	if err != nil {
		return false, err
	}
	normalizedRemoteState, err := r.normalizeJSON(r.remoteState)
	if err != nil {
		return false, err
	}
	return !reflect.DeepEqual(normalizedLocalState, normalizedRemoteState), nil
}

func (r Resource) normalizeJSON(data []byte) ([]byte, error) {
	var localJSON map[string]interface{}
	if err := json.Unmarshal(data, &localJSON); err != nil {
		return nil, err
	}
	if _, exists := localJSON["updateTs"]; exists {
		delete(localJSON, "updateTs")
	}

	localJSONByte, err := json.Marshal(localJSON)
	if err != nil {
		return nil, err
	}

	return localJSONByte, nil
}

// Resourcer interface for Spinnaker resource
type Resourcer interface {
	// IsChanged is used to compare local and remmote state
	IsChanged() (bool, error)
	// SaveRemoteState is used to save state remotely
	SaveRemoteState() error
	// SaveLocalState is used to save state localy
	SaveLocalState() ([]byte, error)
}
