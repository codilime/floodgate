package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/gateclient"
)

// Pipeline object
type Pipeline struct {
	*Resource
	appName string
}

// Init initialize pipeline
func (p *Pipeline) Init(name string, appName string, api *gateclient.GateapiClient, localData []byte) error {
	p.Resource = &Resource{
		name:         name,
		localState:   localData,
		spinnakerAPI: api,
	}
	p.appName = appName
	err := p.loadRemoteState()
	if err != nil {
		return err
	}
	return nil
}

// LoadRemoteState get remote resource
func (p *Pipeline) loadRemoteState() error {
	successPayload, resp, err := p.spinnakerAPI.ApplicationControllerApi.GetPipelineConfigUsingGET(p.spinnakerAPI.Context, p.appName, p.name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline in pipeline %s with name %s, status code: %d", p.appName, p.name, resp.StatusCode)
	}

	jsonPayload, err := json.Marshal(successPayload)
	if err != nil {
		return err
	}
	p.remoteState = jsonPayload

	return nil
}

// SaveRemoteState is used to save state remotely
func (p Pipeline) SaveRemoteState() error {
	var jsonPayload interface{}
	err := json.Unmarshal(p.localState, &jsonPayload)
	if err != nil {
		return err
	}

	saveResp, err := p.spinnakerAPI.PipelineControllerApi.SavePipelineUsingPOST(p.spinnakerAPI.Context, jsonPayload)
	if err != nil {
		return err
	}
	if saveResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d", saveResp.StatusCode)
	}

	return nil
}
