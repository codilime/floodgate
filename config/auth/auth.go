package auth

import "golang.org/x/oauth2"

type Config struct {
	Basic  Basic  `json:"basic"`
	OAuth2 OAuth2 `json:"oauth2"`
	X509   X509   `json:"x509"`
}

type Basic struct {
	Enabled  bool   `json:"enabled"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type OAuth2 struct {
	Enabled      bool         `json:"enabled"`
	TokenURL     string       `json:"tokenUrl"`
	AuthURL      string       `json:"authUrl"`
	ClientID     string       `json:"clientId"`
	ClientSecret string       `json:"clientSecret"`
	Scopes       []string     `json:"scopes"`
	CachedToken  oauth2.Token `json:"cachedToken,omitempty"`
}

type X509 struct {
	Enabled  bool   `json:"enabled"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

// IsAuthValid is used to check if is only one auth method selected
func (config *Config) IsValid() bool {
	if config.Basic.Enabled && config.OAuth2.Enabled {
		return false
	}

	if config.Basic.Enabled && config.X509.Enabled {
		return false
	}

	if config.OAuth2.Enabled && config.X509.Enabled {
		return false
	}

	if !config.Basic.Enabled && !config.OAuth2.Enabled && !config.X509.Enabled {
		return false
	}

	return true
}
