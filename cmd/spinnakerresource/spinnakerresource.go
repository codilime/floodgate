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
	var (
		localJSON, remoteJSON map[string]interface{}
	)

	if err := json.Unmarshal(r.localState, &localJSON); err != nil {
		return false, err
	}
	if err := json.Unmarshal(r.remoteState, &remoteJSON); err != nil {
		return false, err
	}

	for k := range remoteJSON {
		if _, exists := localJSON[k]; !exists {
			delete(remoteJSON, k)
		}
	}

	return !reflect.DeepEqual(localJSON, remoteJSON), nil
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
