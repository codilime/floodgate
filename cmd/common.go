package cmd

import (
	c "github.com/codilime/floodgate/config"
	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	rh "github.com/codilime/floodgate/resourcehandler"
)

func getResourceHandler(configPath string) (*rh.ResourceHandler, error) {
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	client := gateclient.NewGateapiClient(config)
	p := parser.CreateParser(config.Libraries)
	if err := p.LoadObjectsFromDirectories(config.Resources); err != nil {
		return nil, err
	}
	resourceHandler := &rh.ResourceHandler{}
	if err := resourceHandler.Init(client, &p.Resources); err != nil {
		return nil, err
	}
	return resourceHandler, nil
}
