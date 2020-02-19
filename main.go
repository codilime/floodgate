package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	content, err := ioutil.ReadFile("/tmp/pipeline.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	pipeline := new(spr.Pipeline)
	pipeline.Init("deploy-nginx", "nginx", client, content)
	needToSave, err := pipeline.IsChanged()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if needToSave {
		fmt.Print("Saving local state to Spinnaker\n")
		err := pipeline.SaveRemoteState()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	} else {
		fmt.Print("No need to save")
	}
	fmt.Printf("%T\n", pipeline)

}
