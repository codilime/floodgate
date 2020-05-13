package resourcemanager

import (
	"encoding/json"
	"fmt"
	swagger "github.com/codilime/floodgate/gateapi"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	c "github.com/codilime/floodgate/config"
	gc "github.com/codilime/floodgate/gateclient"
	p "github.com/codilime/floodgate/parser"
	fl "github.com/codilime/floodgate/parser/fileloader"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/codilime/floodgate/util"
)

// ResourceManager stores Spinnaker resources and has methods for access, syncing etc.
type ResourceManager struct {
	resources spr.SpinnakerResources
	client    *gc.GateapiClient
}

// Init initializes ResourceManager
func (rm *ResourceManager) Init(config *c.Config, customOptions ...Option) error {
	options := Options{}
	for _, option := range customOptions {
		option(&options)
	}
	if len(options.fileLoaders) == 0 {
		options.fileLoaders = []fl.FileLoader{
			fl.NewJSONLoader(),
			fl.NewJsonnetLoader(config.Libraries...),
			fl.NewYAMLLoader(),
		}
	}
	rm.client = gc.NewGateapiClient(config)
	user, _, err := rm.client.AuthControllerApi.LoggedOutUsingGET(rm.client.Context)
	if err != nil {
		return err
	}
	if user == "" {
		return fmt.Errorf("authenticating with Spinnaker failed. check if credentials are valid")
	}
	parser, err := p.NewParser(options.fileLoaders...)
	if err != nil {
		return err
	}
	resourceData, err := parser.ParseDirectories(config.Resources)
	if err != nil {
		return err
	}
	rm.createResourcesFromData(resourceData)
	return nil
}

// GetChanges get resources' changes
func (rm ResourceManager) GetChanges() (changes []ResourceChange) {
	for _, application := range rm.resources.Applications {
		var change string
		changed, err := application.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if changed {
			change = application.GetFullDiff()
			changes = append(changes, ResourceChange{Type: "application", ID: "", Name: application.Name(), Changes: change})
		}
	}
	for _, pipeline := range rm.resources.Pipelines {
		var change string
		changed, err := pipeline.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if changed {
			change = pipeline.GetFullDiff()
			changes = append(changes, ResourceChange{Type: "pipeline", ID: pipeline.ID(), Name: pipeline.Name(), Changes: change})
		}
	}
	for _, pipelineTemplate := range rm.resources.PipelineTemplates {
		var change string
		changed, err := pipelineTemplate.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if changed {
			change = pipelineTemplate.GetFullDiff()
			changes = append(changes, ResourceChange{Type: "pipelinetemplate", ID: pipelineTemplate.ID(), Name: pipelineTemplate.Name(), Changes: change})
		}
	}
	return
}

// SyncResources synchronize resources with Spinnaker
func (rm *ResourceManager) SyncResources() error {
	if err := rm.syncApplications(); err != nil {
		log.Fatal(err)
	}
	if err := rm.syncPipelines(); err != nil {
		log.Fatal(err)
	}
	if err := rm.syncPipelineTemplates(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (rm ResourceManager) syncResource(resource spr.Resourcer) (bool, error) {
	needToSave, err := resource.IsChanged()
	if err != nil {
		return false, err
	}
	if !needToSave {
		return false, nil
	}
	if err := resource.SaveLocalState(rm.client); err != nil {
		return false, err
	}
	return true, nil
}

func (rm ResourceManager) syncApplications() error {
	log.Print("Syncing applications...")
	for _, application := range rm.resources.Applications {
		synced, err := rm.syncResource(application)
		if err != nil {
			rm.printResourceError(err)
			return fmt.Errorf("failed to sync application: \"%s\"", application.Name())
		}
		if !synced {
			log.Printf("Application \"%s\": up to date", application.Name())
		} else {
			log.Printf("Application \"%s\": synchronized", application.Name())
		}
	}
	return nil
}

func (rm ResourceManager) syncPipelines() error {
	log.Print("Syncing pipelines...")
	for _, pipeline := range rm.resources.Pipelines {
		synced, err := rm.syncResource(pipeline)
		if err != nil {
			rm.printResourceError(err)
			return fmt.Errorf("failed to sync pipeline: \"%s\"", pipeline.Name())
		}
		if !synced {
			log.Printf("Pipeline \"%s\": up to date", pipeline.Name())
		} else {
			log.Printf("Pipeline \"%s\": synchronized", pipeline.Name())
		}
	}
	return nil
}

func (rm ResourceManager) syncPipelineTemplates() error {
	log.Print("Syncing pipeline templates...")
	for _, pipelineTemplate := range rm.resources.PipelineTemplates {
		synced, err := rm.syncResource(pipelineTemplate)
		if err != nil {
			rm.printResourceError(err)
			return fmt.Errorf("failed to sync pipeline template: \"%s\"", pipelineTemplate.Name())
		}
		if !synced {
			log.Printf("Pipeline template \"%s\": up to date", pipelineTemplate.Name())
		} else {
			log.Printf("Pipeline template \"%s\": synchronized", pipelineTemplate.Name())
		}
	}
	return nil
}

// GetAllApplicationsRemoteState returns a concatenated string of applications JSONs.
func (rm *ResourceManager) GetAllApplicationsRemoteState() (state string) {
	for _, application := range rm.resources.Applications {
		state += string(application.GetRemoteState())
	}
	return
}

// GetAllPipelinesRemoteState returns a concatenated string of pipelines JSONs.
func (rm *ResourceManager) GetAllPipelinesRemoteState() (state string) {
	for _, pipeline := range rm.resources.Pipelines {
		state += string(pipeline.GetRemoteState())
	}
	return
}

// GetAllPipelineTemplatesRemoteState returns a concatenated string of pipeline templates JSONs.
func (rm *ResourceManager) GetAllPipelineTemplatesRemoteState() (state string) {
	for _, pipelineTemplate := range rm.resources.Applications {
		state += string(pipelineTemplate.GetRemoteState())
	}
	return
}

func (rm *ResourceManager) createResourcesFromData(resourceData *p.ParsedResourceData) error {
	for _, localData := range resourceData.Applications {
		application := &spr.Application{}
		if err := application.Init(rm.client, localData); err != nil {
			return err
		}
		rm.resources.Applications = append(rm.resources.Applications, application)
	}
	for _, localData := range resourceData.Pipelines {
		pipeline := &spr.Pipeline{}
		if err := pipeline.Init(rm.client, localData); err != nil {
			return err
		}
		rm.resources.Pipelines = append(rm.resources.Pipelines, pipeline)
	}
	for _, localData := range resourceData.PipelineTemplates {
		pipelineTemplate := &spr.PipelineTemplate{}
		if err := pipelineTemplate.Init(rm.client, localData); err != nil {
			return err
		}
		rm.resources.PipelineTemplates = append(rm.resources.PipelineTemplates, pipelineTemplate)
	}
	return nil
}

// SaveResources save resources
func (rm ResourceManager) SaveResources(dirPath string) error {
	applicationsDir := filepath.Join(dirPath, "applications")
	pipelinesDir := filepath.Join(dirPath, "pipelines")
	pipelineTemplatesDir := filepath.Join(dirPath, "pipelinetemplates")
	util.CreateDirs(applicationsDir, pipelinesDir, pipelineTemplatesDir)
	jsonFileExt := ".json"
	for _, application := range rm.resources.Applications {
		filePath := filepath.Join(applicationsDir, application.Name()+jsonFileExt)
		rm.saveResource(filePath, application)
	}
	for _, pipeline := range rm.resources.Pipelines {
		filePath := filepath.Join(pipelinesDir, pipeline.ID()+jsonFileExt)
		rm.saveResource(filePath, pipeline)
	}
	for _, pipelineTemplate := range rm.resources.PipelineTemplates {
		filePath := filepath.Join(pipelineTemplatesDir, pipelineTemplate.ID()+jsonFileExt)
		rm.saveResource(filePath, pipelineTemplate)
	}
	return nil
}

func (rm ResourceManager) saveResource(filePath string, resource spr.Resourcer) error {
	localData := resource.LocalState()
	return rm.saveResourceData(filePath, localData)
}

func (rm ResourceManager) saveResourceData(filePath string, resourceData []byte) error {
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

func (rm ResourceManager) printResourceError(err error) {
	if swaggerError, ok := err.(swagger.GenericSwaggerError); ok {
		var resErr map[string]interface{}
		err = json.Unmarshal(swaggerError.Body(), &resErr)
		if err != nil {
			log.Errorf("%v", err)
		}

		log.Errorf("Status: %v, Error: %s, Message: %s", resErr["status"], resErr["error"], resErr["message"])
	} else {
		log.Errorf("%v", err)
	}
}
