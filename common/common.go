package common

import (
	c "github.com/codilime/floodgate/config"
	"github.com/codilime/floodgate/gateclient"
	"github.com/codilime/floodgate/parser"
	rm "github.com/codilime/floodgate/resourcemanager"
)

// GetResourceManager return a pointer to a initialized ResourceManager instance.
func GetResourceManager(configPath string) (*rm.ResourceManager, error) {
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	client := gateclient.NewGateapiClient(config)
	p := parser.CreateParser(config.Libraries)
	if err := p.LoadObjectsFromDirectories(config.Resources); err != nil {
		return nil, err
	}
	resourceManager := &rm.ResourceManager{}
	if err := resourceManager.Init(client, &p.Resources); err != nil {
		return nil, err
	}
	return resourceManager, nil
}
