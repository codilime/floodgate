package main

import (
	"encoding/json"
	"log"

	"github.com/codilime/floodgate/cmd/cli"
	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/cmd/parser"
	spr "github.com/codilime/floodgate/cmd/spinnakerresource"
)

func main() {
	floodgateConfig, _ := cli.LoadConfig("config.yaml")

	p := parser.CreateParser(floodgateConfig.Libraries)

	_ = p.LoadObjectsFromDirectories(floodgateConfig.Resources)
	log.Print("resources: ", p.Resources)

	client := gateclient.NewGateapiClient(floodgateConfig)
	for _, pipelineTemplate := range p.Resources.PipelineTemplates {
		pipelineTemplateJSON, err := json.Marshal(pipelineTemplate)
		if err != nil {
			log.Fatal(err)
		}
		id := pipelineTemplate["id"].(string)
		metadata := pipelineTemplate["metadata"].(map[string]interface{})
		name := metadata["name"].(string)
		newPipelineTemplate := &spr.PipelineTemplate{}
		err = newPipelineTemplate.Init(id, name, client, pipelineTemplateJSON)
		if err != nil {
			log.Fatal(err)
		}
		needToSave, err := newPipelineTemplate.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if needToSave {
			log.Printf("Saving local state of pipeline template with id %s to Spinnaker\n", id)
			err := newPipelineTemplate.SaveRemoteState()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("Pipeline template with id '%s': no need to save", id)
		}
	}
	for _, pipeline := range p.Resources.Pipelines {
		pipelineJSON, err := json.Marshal(pipeline)
		if err != nil {
			log.Fatal(err)
		}
		pipelineName := pipeline["name"].(string)
		pipelineApp := pipeline["application"].(string)
		newPipeline := &spr.Pipeline{}
		newPipeline.Init(pipelineName, pipelineApp, client, pipelineJSON)

		needToSave, err := newPipeline.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if needToSave {
			log.Print("Saving local state to Spinnaker\n")
			err := newPipeline.SaveRemoteState()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Print("No need to save")
		}
	}
}
