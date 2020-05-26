package config

import (
	"fmt"
	"golang.org/x/oauth2"
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
	Endpoint string `json:"endpoint"`
	Insecure bool   `json:"insecure"`
	// TODO(wurbanski): use other auths than basic
	Auth struct {
		Basic struct {
			Enabled  bool   `json:"enabled"`
			User     string `json:"user"`
			Password string `json:"password"`
		} `json:"basic"`

		OAuth2 struct {
			Enabled      bool         `json:"enabled"`
			TokenURL     string       `json:"tokenUrl"`
			AuthURL      string       `json:"authUrl"`
			ClientID     string       `json:"clientId"`
			ClientSecret string       `json:"clientSecret"`
			Scopes       []string     `json:"scopes"`
			CachedToken  oauth2.Token `json:"cachedToken,omitempty"`
		} `json:"oauth2"`

		X509 struct {
			Enabled  bool   `json:"enabled"`
			CertPath string `json:"certPath"`
			KeyPath  string `json:"keyPath"`
		} `json:"x509"`
	} `json:"auth"`
	Libraries []string `json:"libraries"`
	Resources []string `json:"resources"`
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
