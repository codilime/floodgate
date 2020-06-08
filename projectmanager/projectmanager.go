package projectmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/codilime/floodgate/config"
	gc "github.com/codilime/floodgate/gateclient"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/codilime/floodgate/util"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// ProjectManager stores resources and data for specified project and has methods for access, saving.
type ProjectManager struct {
	resources         spr.SpinnakerResources
	client            *gc.GateapiClient
	projectName       string
	pipelineConfigIds []string
}

// Init initializes ProjectManager
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
	var apps []interface{}
	var pipelines []interface{}

	if _, exists := payload["config"]; exists {
		config, ok := payload["config"].(map[string]interface{})
		if ok {
			projectConfig = config
		}
	}

	if _, exists := projectConfig["applications"]; exists {
		applications, ok := projectConfig["applications"].([]interface{})
		if ok {
			apps = applications
		}

	}

	if _, exists := projectConfig["pipelineConfigs"]; exists {
		pc, ok := projectConfig["pipelineConfigs"].([]interface{})
		if ok {
			pipelines = pc
		}
	}

	for _, appName := range apps {
		appName, ok := appName.(string)
		if ok {
			app := spr.Application{}
			err := app.LoadRemoteStateByName(pm.client, appName)
			if err != nil {
				return err
			}

			pm.resources.Applications = append(pm.resources.Applications, &app)
		}
	}

	for _, pipelineConfig := range pipelines {
		if pc, ok := pipelineConfig.(map[string]interface{}); ok {
			if configID, ok := pc["pipelineConfigId"].(string); ok {
				pm.pipelineConfigIds = append(pm.pipelineConfigIds, configID)
			}
		}
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
			pipelineConfig, ok := pipeline.(map[string]interface{})
			if !ok {
				return errors.New("malformed data from spinnaker")
			}

			for _, id := range pm.pipelineConfigIds {
				pipelineConfigID, okID := pipelineConfig["id"].(string)
				pipelineConfigName, okName := pipelineConfig["name"].(string)
				if okID && okName && pipelineConfigID == id {
					p := &spr.Pipeline{}
					err := p.LoadRemoteStateByName(pm.client, id, app.Name(), pipelineConfigName)
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
