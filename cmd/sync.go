package cmd

import (
	"fmt"
	"io"
	"log"

	c "github.com/codilime/floodgate/config"
	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	spr "github.com/codilime/floodgate/spinnakerresource"
	"github.com/spf13/cobra"
)

// syncOptions store synchronize command options
type syncOptions struct {
	dryRun bool
}

// NewSyncCmd create synchronize command
func NewSyncCmd(out io.Writer) *cobra.Command {
	options := syncOptions{}
	cmd := &cobra.Command{
		Use:     "synchronize",
		Aliases: []string{"sync"},
		Short:   "Synchronize resources to Spinnaker",
		Long:    "Synchronize resources to Spinnaker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(cmd, options)
		},
	}
	cmd.Flags().BoolVarP(&options.dryRun, "dry-run", "d", false, "process resources and preview the result but don't sync with Spinnaker")
	return cmd
}

func runSync(cmd *cobra.Command, options syncOptions) error {
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	// TODO(mlembke): check if there's no error with client
	client := gateclient.NewGateapiClient(config)
	p := parser.CreateParser(config.Libraries)
	if err := p.LoadObjectsFromDirectories(config.Resources); err != nil {
		return err
	}
	sync := &sync{parser: p, client: client}
	if err := sync.syncResources(); err != nil {
		return err
	}
	return nil
}

type sync struct {
	parser *parser.Parser
	client *gateclient.GateapiClient
}

func (s *sync) syncResources() error {
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

func (s sync) syncResource(resource spr.Resourcer, localData map[string]interface{}) (bool, error) {
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

func (s sync) syncApplications() error {
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

func (s sync) syncPipelines() error {
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

func (s sync) syncPipelineTemplates() error {
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
