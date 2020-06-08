package projectmanager

import (
	"encoding/json"
	"fmt"
	c "github.com/codilime/floodgate/config"
	gc "github.com/codilime/floodgate/gateclient"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/codilime/floodgate/util"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type ProjectManager struct {
	resources         spr.SpinnakerResources
	client            *gc.GateapiClient
	projectName       string
	pipelineConfigIds []string
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

	if err := pm.getApplications(); err != nil {
		return err
	}

	if err := pm.getPipelines(); err != nil {
		return err
	}

	return nil
}

func (pm *ProjectManager) getApplications() error {
	log.Print("Loading project applications...")

	payload, _, err := pm.client.ProjectControllerApi.GetUsingGET1(pm.client.Context, pm.projectName)
	if err != nil {
		return err
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

	for _, pipelineConfig := range projectConfig["pipelineConfigs"].([]interface{}) {
		pc := pipelineConfig.(map[string]interface{})
		pm.pipelineConfigIds = append(pm.pipelineConfigIds, pc["pipelineConfigId"].(string))
	}

	return nil
}

func (pm *ProjectManager) getPipelines() error {
	log.Print("Loading project pipelines...")

	for _, app := range pm.resources.Applications {
		payload, _, err := pm.client.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(pm.client.Context, app.Name())
		if err != nil {
			return err
		}

		for _, pipeline := range payload {
			pipelineConfig := pipeline.(map[string]interface{})
			for _, id := range pm.pipelineConfigIds {
				if pipelineConfig["id"].(string) == id {
					p := &spr.Pipeline{}
					err := p.LoadRemoteStateByName(pm.client, id, app.Name(), pipelineConfig["name"].(string))
					if err != nil {
						return err
					}
					pm.resources.Pipelines = append(pm.resources.Pipelines, p)
				}
			}
		}
	}
	return nil
}

// SaveResources save resources
func (pm ProjectManager) SaveResources(dirPath string) error {
	applicationsDir := filepath.Join(dirPath, pm.projectName, "applications")
	pipelinesDir := filepath.Join(dirPath, pm.projectName, "pipelines")
	util.CreateDirs(applicationsDir, pipelinesDir)
	jsonFileExt := ".json"
	for _, application := range pm.resources.Applications {
		filePath := filepath.Join(applicationsDir, application.Name()+jsonFileExt)
		pm.saveResource(filePath, application)
	}
	for _, pipeline := range pm.resources.Pipelines {
		filePath := filepath.Join(pipelinesDir, pipeline.ID()+jsonFileExt)
		pm.saveResource(filePath, pipeline)
	}
	return nil
}

func (pm ProjectManager) saveResource(filePath string, resource spr.Resourcer) error {
	remoteData := resource.RemoteState()
	return pm.saveResourceData(filePath, remoteData)
}

func (pm ProjectManager) saveResourceData(filePath string, resourceData []byte) error {
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	var obj map[string]interface{}
	json.Unmarshal(resourceData, &obj)
	if err := encoder.Encode(obj); err != nil {
		return err
	}
	return nil
}
