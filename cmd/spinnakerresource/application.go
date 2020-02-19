package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/gateclient"
	gateapi "github.com/codilime/floodgate/gateapi"
)

type Application Resource

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
			a.remoteState, _ = json.Marshal(payload)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application %s, status code: %d", a.name, resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

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
	_, _, err := a.spinnakerAPI.TaskControllerApi.TaskUsingPOST1(a.spinnakerAPI.Context, createApplicationTask)
	if err != nil {
		return err
	}

	// TODO(wurbanski): Check if the application was actually created using TaskController
	return nil
}
