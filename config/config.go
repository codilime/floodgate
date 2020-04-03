package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/ghodss/yaml"
)

// Config is the default configuration for the app
type Config struct {
	Endpoint string `yaml:"endpoint"`
	Insecure bool   `yaml:"insecure"`
	// TODO(wurbanski): use other auths than basic
	Auth struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
	Libraries []string `yaml:"libraries"`
	Resources []string `yaml:"resources"`
}

// LoadConfig function is used to load configuration from file
func LoadConfig(locations ...string) (*Config, error) {
	var location string
	if len(locations) == 0 {
		return nil, fmt.Errorf("no config file provided")
	}
	location = locations[0]
	if location == "" {
		userHome := ""
		usr, err := user.Current()
		if err != nil {
			// Fallback by trying to read $HOME
			userHome = os.Getenv("HOME")
			if userHome != "" {
				err = nil
			} else {
				return nil, fmt.Errorf("failed to read current user from environment: %w", err)
			}
		} else {
			userHome = usr.HomeDir
		}
		location = filepath.Join(userHome, ".config", "floodgate", "config.yaml")
	}

	conf := &Config{}

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
