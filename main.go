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

	err := p.LoadObjectsFromDirectories(floodgateConfig.Resources)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("resources: ", p.Resources)

	client := gateclient.NewGateapiClient(floodgateConfig)
	for _, pipelineTemplate := range p.Resources.PipelineTemplates {
		newPipelineTemplate := &spr.PipelineTemplate{}
		err = newPipelineTemplate.Init(client, pipelineTemplate)
		if err != nil {
			log.Fatalf("Encountered an error while processing pipeline template %v: %v", pipelineTemplate, err)
		}
		needToSave, err := newPipelineTemplate.IsChanged()
		if err != nil {
			log.Fatal(err)
		}
		if needToSave {
			log.Printf("Saving local state of pipeline template %v to Spinnaker\n", pipelineTemplate)
			err := newPipelineTemplate.SaveRemoteState()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("No need to save pipeline template %v", pipelineTemplate)
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
