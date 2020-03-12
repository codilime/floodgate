package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/util"

	"github.com/codilime/floodgate/cmd/gateclient"
)

// PipelineTemplate object
type PipelineTemplate struct {
	*Resource
	id string
}

// Init initialize pipeline template
func (pt *PipelineTemplate) Init(api *gateclient.GateapiClient, localData map[string]interface{}) error {
	if err := pt.validate(localData); err != nil {
		return err
	}
	id := localData["id"].(string)
	localState, err := json.Marshal(localData)
	if err != nil {
		return err
	}
	pt.Resource = &Resource{
		localState:   localState,
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
			return fmt.Errorf("failed to save pipeline template, status code: %d", resp.StatusCode)
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
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
