package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/gateclient"
	gateapi "github.com/codilime/floodgate/gateapi"
)

// Application object
type Application struct {
	*Resource
}

// Init function for Application resource
func (a *Application) Init(name string, api *gateclient.GateapiClient, localdata []byte) {
	a.name = name
	a.spinnakerAPI = api
	a.loadRemoteState()
}

func (a *Application) loadRemoteState() error {
	var optionals gateapi.GetApplicationUsingGETOpts
	payload, resp, err := a.spinnakerAPI.ApplicationControllerApi.GetApplicationUsingGET(a.spinnakerAPI.Context, a.name, &optionals)
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			a.remoteState = []byte("{}")
			return nil
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application %s, status code: %d", a.name, resp.StatusCode)
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		a.remoteState = jsonPayload
	}

	if err != nil {
		return err
	}

	return nil
}

// SaveRemoteState function for saving object to Spinnaker
func (a Application) SaveRemoteState() error {
	var app map[string]interface{}
	json.Unmarshal(a.localState, &app)
	createApplicationTask := map[string]interface{}{
		"job": map[string]interface{}{
			"type":        "createApplication",
			"application": "app",
		},
		"application": a.name,
		"description": "Creating application",
	}
	task, _, err := a.spinnakerAPI.TaskControllerApi.TaskUsingPOST1(a.spinnakerAPI.Context, createApplicationTask)
	if err != nil {
		return err
	}

	// TODO(wurbanski): Check if the application was actually created using TaskController
	if err := a.spinnakerAPI.WaitForSuccessfulTask(task, 5); err != nil {
		return err
	}
	return nil
}
