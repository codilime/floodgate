package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/gateclient"

	"github.com/nsf/jsondiff"
)

// Pipeline object
type Pipeline struct {
	name, id, appName       string
	localState, remoteState []byte
}

// CreatePipeline is used to create new Resource object
func CreatePipeline(name string, id string, appName string, api *gateclient.GateapiClient, localData []byte) Resource {
	ppln := Pipeline{name, id, appName, localData, []byte("{}")}
	ppln.loadRemoteState(api)

	return ppln
}

// IsChanged is used to compare local and remmote state
func (p Pipeline) IsChanged() (bool, error) {
	var options jsondiff.Options

	diffType, _ := jsondiff.Compare(p.localState, p.remoteState, &options)

	if diffType.String() != "FullMatch" {
		return true, nil
	}

	return false, nil
}

// LoadRemoteState get remote resource
func (p *Pipeline) loadRemoteState(api *gateclient.GateapiClient) error {
	successPayload, resp, err := api.ApplicationControllerApi.GetPipelineConfigUsingGET(api.Context, p.appName, p.name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline in pipeline %s with name %s, status code: %d", p.appName, p.name, resp.StatusCode)
	}
	jsonPayload, err := json.MarshalIndent(successPayload, "", " ")
	if err != nil {
		return err
	}
	p.remoteState = jsonPayload

	return nil
}

// SaveRemoteState is used to save state remotely
func (p Pipeline) SaveRemoteState(api *gateclient.GateapiClient) error {
	var jsonPayload interface{}
	err := json.Unmarshal(p.localState, &jsonPayload)
	if err != nil {
		return err
	}

	saveResp, err := api.PipelineControllerApi.SavePipelineUsingPOST(api.Context, jsonPayload)
	if err != nil {
		return err
	}
	if saveResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d", saveResp.StatusCode)
	}

	return nil
}

// SaveLocalState is used to save state localy
func (p Pipeline) SaveLocalState() ([]byte, error) {
	return p.localState, nil
}
