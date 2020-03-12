package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/cmd/util"
	gateapi "github.com/codilime/floodgate/gateapi"
)

// Application object
type Application struct {
	*Resource
	name string
}

// Init function for Application resource
func (a *Application) Init(api *gateclient.GateapiClient, localData map[string]interface{}) error {
	if err := a.validate(localData); err != nil {
		return err
	}
	name := localData["name"].(string)
	localState, err := json.Marshal(localData)
	if err != nil {
		return err
	}
	a.Resource = &Resource{
		localState:   localState,
		spinnakerAPI: api,
	}
	a.name = name
	if err := a.loadRemoteState(); err != nil {
		return err
	}
	return nil
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

func (a Application) validate(localData map[string]interface{}) error {
	errors := []error{}
	if err := util.AssertMapKeyIsString(localData, "name", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsString(localData, "email", true); err != nil {
		errors = append(errors, err)
	}
	return util.CombineErrors(errors)
}
