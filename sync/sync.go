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
	client *gateclient.GateapiClient
	parser *parser.Parser
}

// Init initialize sync
func (s *Sync) Init(client *gateclient.GateapiClient, parser *parser.Parser) {
	s.client = client
	s.parser = parser
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

func (s Sync) syncResource(resource spr.Resourcer, localData map[string]interface{}) (bool, error) {
	if err := resource.Init(s.client, localData); err != nil {
		return false, err
	}
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
	for _, applicationData := range s.parser.Resources.Applications {
		application := &spr.Application{}
		synced, err := s.syncResource(application, applicationData)
		if err != nil {
			return fmt.Errorf("failed to sync application: %v", err)
		}
		if !synced {
			log.Printf("No need to save application %v", applicationData)
		}
	}
	return nil
}

func (s Sync) syncPipelines() error {
	log.Print("Syncing pipelines")
	for _, pipelineData := range s.parser.Resources.Pipelines {
		pipeline := &spr.Pipeline{}
		synced, err := s.syncResource(pipeline, pipelineData)
		if err != nil {
			return err
		}
		if !synced {
			log.Printf("No need to save pipeline %v", pipelineData)
		}
	}
	return nil
}

func (s Sync) syncPipelineTemplates() error {
	log.Print("Syncing pipeline templates")
	for _, pipelineTemplateData := range s.parser.Resources.PipelineTemplates {
		pipelineTemplate := &spr.PipelineTemplate{}
		synced, err := s.syncResource(pipelineTemplate, pipelineTemplateData)
		if err != nil {
			return err
		}
		if !synced {
			log.Printf("No need to save pipeline template %v", pipelineTemplateData)
		}
	}
	return nil
}
