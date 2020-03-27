package spinnakerresource

import (
	"encoding/json"
	"fmt"
	"net/http"

	gateapi "github.com/codilime/floodgate/gateapi"
	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/util"
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
		var remoteData map[string]interface{}
		if _, exists := payload["attributes"]; exists {
			attributes, ok := payload["attributes"].(map[string]interface{})
			if ok {
				remoteData = attributes
			}
		}
		if _, exists := payload["clusters"]; exists {
			clusters, ok := payload["clusters"].(map[string]interface{})
			if ok && len(clusters) > 0 {
				remoteData["clusters"] = clusters
			}
		}
		remoteState, err := json.Marshal(remoteData)
		if err != nil {
			return err
		}
		a.remoteState = remoteState
	}

	if err != nil {
		return err
	}

	return nil
}

// LocalState get local state
func (a Application) LocalState() []byte {
	return a.localState
}

// RemoteState get remote state
func (a Application) RemoteState() []byte {
	return a.remoteState
}

// SaveLocalState save local state to Spinnaker
func (a Application) SaveLocalState() error {
	var app map[string]interface{}
	if err := json.Unmarshal(a.localState, &app); err != nil {
		return fmt.Errorf("failed to unmarshal local state")
	}
	createApplicationTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": app}},
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
	// TODO (mlembke): The Name of the application cannot have hyphens, or it will interfere with the naming convention.
	// TODO (mlembke): see https://docs.armory.io/overview/your-first-application/
	errors := []error{}
	if err := util.AssertMapKeyIsString(localData, "name", true); err != nil {
		errors = append(errors, err)
	}
	if err := util.AssertMapKeyIsString(localData, "email", true); err != nil {
		errors = append(errors, err)
	}
	return util.CombineErrors(errors)
}
