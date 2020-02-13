package cli

import (
	"io/ioutil"
	"log"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/config"
	"gopkg.in/yaml.v2"
)

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
	} else {
		log.Printf("Config: %v\n", conf)
	}
	return conf, nil
}
