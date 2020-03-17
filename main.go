package main

import (
	"log"

	"github.com/codilime/floodgate/cmd/cli"
	"github.com/codilime/floodgate/cmd/gateclient"
	"github.com/codilime/floodgate/cmd/parser"
	"github.com/codilime/floodgate/cmd/sync"
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
	sync := &sync.Sync{}
	sync.Init(p, client)
	if err := sync.Sync(); err != nil {
		log.Fatal(err)
	}
}
