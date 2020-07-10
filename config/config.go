package config

import (
	"errors"
	"fmt"
	"github.com/codilime/floodgate/config/auth"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/ghodss/yaml"
)

var location string

// Config is the default configuration for the app
type Config struct {
	Endpoint  string      `json:"endpoint"`
	Insecure  bool        `json:"insecure"`
	Auth      auth.Config `json:"auth"`
	Libraries []string    `json:"libraries"`
	Resources []string    `json:"resources"`
}

// LoadConfig function is used to load configuration from file
func LoadConfig(locations ...string) (*Config, error) {
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
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return nil, err
	}

	if !conf.Auth.IsValid() {
		return nil, errors.New("incorrect auth configuration")
	}

	return conf, nil
}

// Merge method is used to override current config with new one
func (c *Config) Merge(cfg Config) {
	if cfg.Endpoint != "" {
		c.Endpoint = cfg.Endpoint
	}
}

// SaveConfig function is used to save configuration file
func SaveConfig(config *Config) error {
	configFile, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = ioutil.WriteFile(location, configFile, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
