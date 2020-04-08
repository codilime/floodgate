package resourcehandler

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	spr "github.com/codilime/floodgate/spinnakerresource"
)

// SpinnakerResources Spinnaker resources collection
type SpinnakerResources struct {
	Applications      []*spr.Application
	Pipelines         []*spr.Pipeline
	PipelineTemplates []*spr.PipelineTemplate
}

// ResourceChange store resource change
type ResourceChange struct {
	Type    string
	ID      string
	Name    string
	Changes string
}

// ResourceHandler synchornize resources with Spinnaker
type ResourceHandler struct {
	resources         SpinnakerResources
	desyncedResources SpinnakerResources
}

// Init initialize sync
func (rh *ResourceHandler) Init(client *gateclient.GateapiClient, resourceData *parser.ResourceData) error {
	for _, localData := range resourceData.Applications {
		application := &spr.Application{}
		if err := application.Init(client, localData); err != nil {
			return err
		}
		rh.resources.Applications = append(rh.resources.Applications, application)
		changed, err := application.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			rh.desyncedResources.Applications = append(rh.desyncedResources.Applications, application)
		}
	}
	for _, localData := range resourceData.Pipelines {
		pipeline := &spr.Pipeline{}
		if err := pipeline.Init(client, localData); err != nil {
			return err
		}
		rh.resources.Pipelines = append(rh.resources.Pipelines, pipeline)
		changed, err := pipeline.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			rh.desyncedResources.Pipelines = append(rh.desyncedResources.Pipelines, pipeline)
		}
	}
	for _, localData := range resourceData.PipelineTemplates {
		pipelineTemplate := &spr.PipelineTemplate{}
		if err := pipelineTemplate.Init(client, localData); err != nil {
			return err
		}
		rh.resources.PipelineTemplates = append(rh.resources.PipelineTemplates, pipelineTemplate)
		changed, err := pipelineTemplate.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			rh.desyncedResources.PipelineTemplates = append(rh.desyncedResources.PipelineTemplates, pipelineTemplate)
		}
	}
	return nil
}

// GetChanges get resources' changes
func (rh ResourceHandler) GetChanges() (changes []ResourceChange) {
	for _, application := range rh.resources.Applications {
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
	for _, pipeline := range rh.resources.Pipelines {
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
	for _, pipelineTemplate := range rh.resources.PipelineTemplates {
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
func (rh *ResourceHandler) SyncResources() error {
	if err := rh.syncApplications(); err != nil {
		log.Fatal(err)
	}
	if err := rh.syncPipelines(); err != nil {
		log.Fatal(err)
	}
	if err := rh.syncPipelineTemplates(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (rh ResourceHandler) syncResource(resource spr.Resourcer) (bool, error) {
	needToSave, err := resource.IsChanged()
	if err != nil {
		return false, err
	}
	if !needToSave {
		return false, nil
	}
	if err := resource.SaveLocalState(); err != nil {
		return false, err
	}
	return true, nil
}

func (rh ResourceHandler) syncApplications() error {
	log.Print("Syncing applications")
	for _, application := range rh.resources.Applications {
		synced, err := rh.syncResource(application)
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

func (rh ResourceHandler) syncPipelines() error {
	log.Print("Syncing pipelines")
	for _, pipeline := range rh.resources.Pipelines {
		synced, err := rh.syncResource(pipeline)
		if err != nil {
			return fmt.Errorf("failed to sync pipeline: %v", pipeline)
		}
		if !synced {
			log.Printf("No need to save pipeline %v", pipeline)
		}
	}
	return nil
}

func (rh ResourceHandler) syncPipelineTemplates() error {
	log.Print("Syncing pipeline templates")
	for _, pipelineTemplate := range rh.resources.PipelineTemplates {
		synced, err := rh.syncResource(pipelineTemplate)
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
func (rh *ResourceHandler) GetAllApplicationsRemoteState() (state string) {
	for _, application := range rh.resources.Applications {
		state += string(application.GetRemoteState())
	}
	return
}

// GetAllPipelinesRemoteState returns a concatenated string of pipelines JSONs.
func (rh *ResourceHandler) GetAllPipelinesRemoteState() (state string) {
	for _, pipeline := range rh.resources.Pipelines {
		state += string(pipeline.GetRemoteState())
	}
	return
}

// GetAllPipelineTemplatesRemoteState returns a concatenated string of pipeline templates JSONs.
func (rh *ResourceHandler) GetAllPipelineTemplatesRemoteState() (state string) {
	for _, pipelineTemplate := range rh.resources.Applications {
		state += string(pipelineTemplate.GetRemoteState())
	}
	return
}
