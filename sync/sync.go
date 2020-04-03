package sync

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

// Sync synchornize resources with Spinnaker
type Sync struct {
	resources         SpinnakerResources
	desyncedResources SpinnakerResources
}

// Init initialize sync
func (s *Sync) Init(client *gateclient.GateapiClient, resourceData *parser.ResourceData) error {
	for _, localData := range resourceData.Applications {
		application := &spr.Application{}
		if err := application.Init(client, localData); err != nil {
			return err
		}
		s.resources.Applications = append(s.resources.Applications, application)
		changed, err := application.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			s.desyncedResources.Applications = append(s.desyncedResources.Applications, application)
		}
	}
	for _, localData := range resourceData.Pipelines {
		pipeline := &spr.Pipeline{}
		if err := pipeline.Init(client, localData); err != nil {
			return err
		}
		s.resources.Pipelines = append(s.resources.Pipelines, pipeline)
		changed, err := pipeline.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			s.desyncedResources.Pipelines = append(s.desyncedResources.Pipelines, pipeline)
		}
	}
	for _, localData := range resourceData.PipelineTemplates {
		pipelineTemplate := &spr.PipelineTemplate{}
		if err := pipelineTemplate.Init(client, localData); err != nil {
			return err
		}
		s.resources.PipelineTemplates = append(s.resources.PipelineTemplates, pipelineTemplate)
		changed, err := pipelineTemplate.IsChanged()
		if err != nil {
			return err
		}
		if changed {
			s.desyncedResources.PipelineTemplates = append(s.desyncedResources.PipelineTemplates, pipelineTemplate)
		}
	}
	return nil
}

// GetChanges get resources' changes
func (s Sync) GetChanges() (changes []ResourceChange) {
	for _, application := range s.resources.Applications {
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
	for _, pipeline := range s.resources.Pipelines {
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
	for _, pipelineTemplate := range s.resources.PipelineTemplates {
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
func (s *Sync) SyncResources() error {
	if err := s.syncApplications(); err != nil {
		log.Fatal(err)
	}
	if err := s.syncPipelines(); err != nil {
		log.Fatal(err)
	}
	if err := s.syncPipelineTemplates(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s Sync) syncResource(resource spr.Resourcer) (bool, error) {
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

func (s Sync) syncApplications() error {
	log.Print("Syncing applications")
	for _, application := range s.resources.Applications {
		synced, err := s.syncResource(application)
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

func (s Sync) syncPipelines() error {
	log.Print("Syncing pipelines")
	for _, pipeline := range s.resources.Pipelines {
		synced, err := s.syncResource(pipeline)
		if err != nil {
			return fmt.Errorf("failed to sync pipeline: %v", pipeline)
		}
		if !synced {
			log.Printf("No need to save pipeline %v", pipeline)
		}
	}
	return nil
}

func (s Sync) syncPipelineTemplates() error {
	log.Print("Syncing pipeline templates")
	for _, pipelineTemplate := range s.resources.PipelineTemplates {
		synced, err := s.syncResource(pipelineTemplate)
		if err != nil {
			return fmt.Errorf("failed to sync pipeline template: %v", pipelineTemplate)
		}
		if !synced {
			log.Printf("No need to save pipeline template %v", pipelineTemplate)
		}
	}
	return nil
}
