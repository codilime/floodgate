package main

import (
	"log"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/cli"
	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/parser"
)

func main() {
	floodgateConfig, _ := cli.LoadConfig("config.yaml")

	p := parser.CreateParser(floodgateConfig.Libraries)

	_ = p.LoadObjectsFromDirectories(floodgateConfig.Resources)
	log.Print("resources: ", p.Resources)

}
