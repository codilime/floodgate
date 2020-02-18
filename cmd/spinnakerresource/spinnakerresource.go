package spinnakerresource

import "cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/gateclient"

// Resourcer interface for Spinnaker resource
type Resourcer interface {
	// IsChanged is used to compare local and remmote state
	IsChanged() (bool, error)
	// SaveRemoteState is used to save state remotely
	SaveRemoteState(api *gateclient.GateapiClient) error
	// SaveLocalState is used to save state localy
	SaveLocalState() ([]byte, error)
}
