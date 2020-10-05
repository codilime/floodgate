package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	gc "github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/util"
)

// Pipeline object
type Pipeline struct {
	*Resource
	name              string
	application       string
	id                string
	templateReference string
}

// Init initialize pipeline
func (p *Pipeline) Init(api *gc.GateapiClient, localData map[string]interface{}) error {
	if err := p.validate(localData); err != nil {
		return err
	}
	name := localData["name"].(string)
	application := localData["application"].(string)
	id := localData["id"].(string)
	reference := p.getTemplateReference(localData)

	localState, err := json.Marshal(localData)
	if err != nil {
		return err
	}
	p.Resource = &Resource{
		localState: localState,
	}
	p.name = name
	p.application = application
	p.id = id
	p.templateReference = reference
	if api != nil {
		if err := p.LoadRemoteState(api); err != nil {
			err := p.SaveLocalState(api)
			return err
		}
	}
	return nil
}

// Name get pipeline name
func (p Pipeline) Name() string {
	return p.name
}

// ID get pipeline id
func (p Pipeline) ID() string {
	return p.id
}

// Application get application name
func (p Pipeline) Application() string {
	return p.application
}

// TemplateReference get template reference
func (p Pipeline) TemplateReference() string {
	return p.templateReference
}

// LoadRemoteStateByName load resource's remote state from Spinnaker by provided name
func (p *Pipeline) LoadRemoteStateByName(spinnakerAPI *gc.GateapiClient, id, application, name string) error {
	p.name = name
	p.application = application
	p.id = id
	p.Resource = &Resource{}

	return p.LoadRemoteState(spinnakerAPI)
}

// LoadRemoteState load resource's remote state from Spinnaker
func (p *Pipeline) LoadRemoteState(spinnakerAPI *gc.GateapiClient) error {
	successPayload, resp, err := spinnakerAPI.ApplicationControllerApi.GetPipelineConfigUsingGET(spinnakerAPI.Context, p.application, p.name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline in pipeline %s with name %s, status code: %d", p.application, p.name, resp.StatusCode)
	}

	jsonPayload, err := json.Marshal(successPayload)
	if err != nil {
		return err
	}
	p.remoteState = jsonPayload

	return nil
}

// LocalState get local state
func (p Pipeline) LocalState() []byte {
	return p.localState
}

// RemoteState get remote state
func (p Pipeline) RemoteState() []byte {
	return p.remoteState
}

// SaveLocalState save local state to Spinnaker
func (p Pipeline) SaveLocalState(spinnakerAPI *gc.GateapiClient) error {
	var jsonPayload interface{}
	err := json.Unmarshal(p.localState, &jsonPayload)
	if err != nil {
		return err
	}

	saveResp, err := spinnakerAPI.PipelineControllerApi.SavePipelineUsingPOST(spinnakerAPI.Context, jsonPayload)
	if err != nil {
		return err
	}
	if saveResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d", saveResp.StatusCode)
	}

	return nil
}

func (p Pipeline) getTemplateReference(localData map[string]interface{}) string {
	err := util.AssertMapKeyIsStringMap(localData, "template", true)
	if err == nil {
		template := localData["template"].(map[string]interface{})
		return template["reference"].(string)
	}
	return ""
}

func (p Pipeline) validate(localData map[string]interface{}) error {
	if err := util.AssertMapKeyIsString(localData, "name", true); err != nil {
		return err
	}
	if err := util.AssertMapKeyIsString(localData, "application", true); err != nil {
		return err
	}
	if err := util.AssertMapKeyIsString(localData, "id", true); err != nil {
		return err
	}
	// template is optional key, but it requires defined schema
	err := util.AssertMapKeyIsStringMap(localData, "template", true)
	if err == nil {
		if err := util.AssertMapKeyIsString(localData, "schema", true); err != nil {
			return err
		}
	}
	return nil
}
