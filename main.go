package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/cli"
	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/gateclient"
	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/parser"
	spr "cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/spinnakerresource"
)

func main() {
	floodgateConfig, _ := cli.LoadConfig()

	p, _ := parser.CreateParser(floodgateConfig.Libraries)

	output, _ := p.LoadDirectories(floodgateConfig.Resources)
	log.Print(output)

	client := gateclient.NewGateapiClient(floodgateConfig)
	content, err := ioutil.ReadFile("/tmp/pipeline.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	var pipeline spr.Resource
	pipeline = spr.CreatePipeline("deploy-nginx", "1bfaa7c1-894c-4adb-9e51-c969bc38c984", "nginx", client, content)
	needToSave, err := pipeline.IsChanged()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if needToSave {
		fmt.Print("Saving local state to Spinnaker\n")
		err := pipeline.SaveRemoteState(client)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	} else {
		fmt.Print("No need to save")
	}
	fmt.Printf("%T\n", pipeline)
}
