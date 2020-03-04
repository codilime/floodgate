package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/gateclient"
)

// PipelineTemplate object
type PipelineTemplate struct {
	*Resource
	id string
}

// Init initialize pipeline template
func (pt *PipelineTemplate) Init(id string, name string, api *gateclient.GateapiClient, localData []byte) error {
	pt.Resource = &Resource{
		name:         name,
		localState:   localData,
		spinnakerAPI: api,
	}
	pt.id = id
	if err := pt.loadRemoteState(); err != nil {
		return err
	}
	return nil
}

// loadRemoteState load resource state from Spinnaker
func (pt *PipelineTemplate) loadRemoteState() error {
	successPayload, resp, err := pt.spinnakerAPI.V2PipelineTemplatesControllerApi.GetUsingGET2(pt.spinnakerAPI.Context, pt.id, nil)
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			pt.remoteState = []byte("{}")
			return nil
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error while getting pipeline template with id %s, status code: %d", pt.id, resp.StatusCode)
		}
		jsonPayload, err := json.Marshal(successPayload)
		if err != nil {
			return err
		}
		pt.remoteState = jsonPayload
	}
	if err != nil {
		return err
	}
	return nil
}

// SaveRemoteState save local state to Spinnaker
func (pt PipelineTemplate) SaveRemoteState() error {
	var localStateJSON interface{}
	err := json.Unmarshal(pt.localState, &localStateJSON)
	if err != nil {
		return err
	}
	var resp *http.Response
	if string(pt.remoteState) == "{}" {
		resp, err = pt.spinnakerAPI.V2PipelineTemplatesControllerApi.CreateUsingPOST1(pt.spinnakerAPI.Context, localStateJSON, nil)
	} else {
		resp, err = pt.spinnakerAPI.V2PipelineTemplatesControllerApi.UpdateUsingPOST1(pt.spinnakerAPI.Context, pt.id, localStateJSON, nil)
	}
	if resp != nil {
		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("Encountered an error saving pipeline, status code: %d", resp.StatusCode)
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
