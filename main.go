package main

import (
	"log"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/cli"
	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/parser"
)

func main() {
	floodgateConfig, _ := cli.LoadConfig()

	p, _ := parser.CreateParser(floodgateConfig.Libraries)

	output, _ := p.LoadDirectories(floodgateConfig.Resources)
	log.Print(output)

}
