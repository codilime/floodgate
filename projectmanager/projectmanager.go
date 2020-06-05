package projectmanager

import (
	"fmt"
	c "github.com/codilime/floodgate/config"
	gc "github.com/codilime/floodgate/gateclient"
	spr "github.com/codilime/floodgate/spinnakerresource"
	log "github.com/sirupsen/logrus"
)

type ProjectManager struct {
	resources   spr.SpinnakerResources
	client      *gc.GateapiClient
	projectName string
}

func (pm *ProjectManager) Init(config *c.Config, projectName string) error {
	pm.projectName = projectName
	pm.client = gc.NewGateapiClient(config)
	user, _, err := pm.client.AuthControllerApi.LoggedOutUsingGET(pm.client.Context)
	if err != nil {
		return err
	}
	if user == "" {
		return fmt.Errorf("authenticating with Spinnaker failed. check if credentials are valid")
	}

	err = pm.GetProject()
	if err != nil {
		return err
	}

	return nil
}

func (pm *ProjectManager) GetProject() error {
	log.Print("Getting project details...")

	payload, _, err := pm.client.ProjectControllerApi.GetUsingGET1(pm.client.Context, pm.projectName)
	if err != nil {
		log.Fatal(err)
	}

	var projectConfig map[string]interface{}
	if _, exists := payload["config"]; exists {
		config, ok := payload["config"].(map[string]interface{})
		if ok {
			projectConfig = config
		}
	}

	for _, appName := range projectConfig["applications"].([]interface{}) {
		app := spr.Application{}
		err := app.LoadRemoteStateByName(pm.client, appName.(string))
		if err != nil {
			return err
		}

		pm.resources.Applications = append(pm.resources.Applications, &app)
	}

	return nil
}
