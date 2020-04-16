package resourcemanager

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	c "github.com/codilime/floodgate/config"
	gc "github.com/codilime/floodgate/gateclient"
	p "github.com/codilime/floodgate/parser"
	spr "github.com/codilime/floodgate/spinnakerresource"
)

// ResourceChange store resource change
type ResourceChange struct {
	Type    string
	ID      string
	Name    string
	Changes string
}

// ResourceManager stores Spinnaker resources and has methods for access, syncing etc.
type ResourceManager struct {
	resources spr.SpinnakerResources
	client    *gc.GateapiClient
}

// Init initialize sync
func (rm *ResourceManager) Init(configPath string) error {
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	rm.client = gc.NewGateapiClient(config)
	parser, err := p.NewResourceParser(config.Libraries...)
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
	log.Print("Syncing applications")
	for _, application := range rm.resources.Applications {
		synced, err := rm.syncResource(application)
		if err != nil {
			return fmt.Errorf("failed to sync application: %v", application)
		}
		if !synced {
			log.Printf("No need to save application %v", application)
		} else {
			log.Printf("Successfully synced application %v", application)
		}
	}
	return nil
}

func (rm ResourceManager) syncPipelines() error {
	log.Print("Syncing pipelines")
	for _, pipeline := range rm.resources.Pipelines {
		synced, err := rm.syncResource(pipeline)
		if err != nil {
			return fmt.Errorf("failed to sync pipeline: %v", pipeline)
		}
		if !synced {
			log.Printf("No need to save pipeline %v", pipeline)
		}
	}
	return nil
}

func (rm ResourceManager) syncPipelineTemplates() error {
	log.Print("Syncing pipeline templates")
	for _, pipelineTemplate := range rm.resources.PipelineTemplates {
		synced, err := rm.syncResource(pipelineTemplate)
		if err != nil {
			return fmt.Errorf("failed to sync pipeline template: %v", pipelineTemplate)
		}
		if !synced {
			log.Printf("No need to save pipeline template %v", pipelineTemplate)
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
