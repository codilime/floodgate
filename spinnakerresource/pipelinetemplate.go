package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/util"

	gc "github.com/codilime/floodgate/gateclient"
)

// PipelineTemplate object
type PipelineTemplate struct {
	*Resource
	id   string
	name string
}

// Init initialize pipeline template
func (pt *PipelineTemplate) Init(api *gc.GateapiClient, localData map[string]interface{}) error {
	if err := pt.validate(localData); err != nil {
		return err
	}
	id := localData["id"].(string)
	name := localData["metadata"].(map[string]interface{})["name"].(string)
	localState, err := json.Marshal(localData)
	if err != nil {
		return err
	}
	pt.Resource = &Resource{
		localState: localState,
	}
	pt.id = id
	pt.name = name
	if api != nil {
		if err := pt.LoadRemoteState(api); err != nil {
			err := pt.SaveLocalState(api)
			return err
		}
	}
	return nil
}

// ID get pipeline template id
func (pt PipelineTemplate) ID() string {
	return pt.id
}

// Name get pipeline template name
func (pt PipelineTemplate) Name() string {
	return pt.name
}

// LoadRemoteState load resource's remote state from Spinnaker
func (pt *PipelineTemplate) LoadRemoteState(spinnakerAPI *gc.GateapiClient) error {
	successPayload, resp, err := spinnakerAPI.V2PipelineTemplatesControllerApi.GetUsingGET2(spinnakerAPI.Context, pt.id, nil)
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			pt.remoteState = []byte("{}")
			return nil
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get pipeline template remote state, status code: %d", resp.StatusCode)
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

// SaveLocalState save local state to Spinnaker
func (pt PipelineTemplate) SaveLocalState(spinnakerAPI *gc.GateapiClient) error {
	var localStateJSON interface{}
	err := json.Unmarshal(pt.localState, &localStateJSON)
	if err != nil {
		return err
	}
	var resp *http.Response
	if string(pt.remoteState) == "{}" {
		resp, err = spinnakerAPI.V2PipelineTemplatesControllerApi.CreateUsingPOST1(spinnakerAPI.Context, localStateJSON, nil)
	} else {
		resp, err = spinnakerAPI.V2PipelineTemplatesControllerApi.UpdateUsingPOST1(spinnakerAPI.Context, pt.id, localStateJSON, nil)
	}
	if resp != nil {
		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("failed to save pipeline template, status code: %d", resp.StatusCode)
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// LocalState get local state
func (pt PipelineTemplate) LocalState() []byte {
	return pt.localState
}

// RemoteState get remote state
func (pt PipelineTemplate) RemoteState() []byte {
	return pt.remoteState
}

// SaveRemoteState save remote state to local storage
func (pt PipelineTemplate) SaveRemoteState() error {
	return fmt.Errorf("Not implemented")
}

func (pt PipelineTemplate) validate(localData map[string]interface{}) error {
	errors := []error{}
	if err := util.AssertMapKeyIsString(localData, "id", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsString(localData, "schema", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsStringMap(localData, "metadata", true); err != nil {
		errors = append(errors, err)
		return util.CombineErrors(errors)
	}
	metadata := localData["metadata"].(map[string]interface{})
	if err := pt.validateMetadata(metadata); err != nil {
		errors = append(errors, err)
	}
	return util.CombineErrors(errors)
}

func (pt PipelineTemplate) validateMetadata(metadata map[string]interface{}) error {
	errors := []error{}
	if err := util.AssertMapKeyIsString(metadata, "name", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsString(metadata, "description", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsString(metadata, "owner", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsInterfaceArray(metadata, "scopes", true); err != nil {
		errors = append(errors, err)
	}
	return util.CombineErrors(errors)
}
