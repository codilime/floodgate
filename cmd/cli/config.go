package cli

import (
	"io/ioutil"
	"log"

	"github.com/codilime/floodgate/config"
	"github.com/ghodss/yaml"
)

// LoadConfig function is used to load configuration from file
func LoadConfig(locations ...string) (*config.Config, error) {
	var location string
	if len(locations) == 0 {
		location = config.DefaultLocation
	} else {
		location = locations[0]
	}

	conf := &config.Config{}

	configFile, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conf, nil

}
