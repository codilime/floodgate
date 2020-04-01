package sync

import (
	"fmt"
	"log"

	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	spr "github.com/codilime/floodgate/spinnakerresource"
)

// Sync synchornize resources with Spinnaker
type Sync struct {
	resources struct {
		Applications      []*spr.Application
		Pipelines         []*spr.Pipeline
		PipelineTemplates []*spr.PipelineTemplate
	}
}

// Init initialize sync
func (s *Sync) Init(client *gateclient.GateapiClient, resourceData *parser.ResourceData) error {
	for _, localData := range resourceData.Applications {
		application := &spr.Application{}
		if err := application.Init(client, localData); err != nil {
			return err
		}
		s.resources.Applications = append(s.resources.Applications, application)
	}
	for _, localData := range resourceData.Pipelines {
		pipeline := &spr.Pipeline{}
		if err := pipeline.Init(client, localData); err != nil {
			return err
		}
		s.resources.Pipelines = append(s.resources.Pipelines, pipeline)
	}
	for _, localData := range resourceData.PipelineTemplates {
		pipelineTemplate := &spr.PipelineTemplate{}
		if err := pipelineTemplate.Init(client, localData); err != nil {
			return err
		}
		s.resources.PipelineTemplates = append(s.resources.PipelineTemplates, pipelineTemplate)

	}
	return nil
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
