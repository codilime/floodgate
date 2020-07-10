package config

import (
	"github.com/codilime/floodgate/config/auth"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

var exampleConfig = &Config{
	Endpoint: "https://127.0.0.1/api/v1",
	Insecure: true,
	Auth: auth.Config{
		Basic: auth.Basic{
			Enabled:  true,
			User:     "admin",
			Password: "VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi",
		},
	},
	Libraries: []string{"sponnet"},
	Resources: []string{"resources/applications", "resources/pipelines", "resources/pipelinetemplates"},
}

func createTempConfigs() (string, string, string, string, error) {
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		return "", "", "", "", err
	}

	valid, err := ioutil.TempFile(dir, "config.yaml")
	if err != nil {
		return "", "", "", "", err
	}
	ioutil.WriteFile(valid.Name(), validConfig, 0644)

	malformed, err := ioutil.TempFile(dir, "okConfig.yaml")
	if err != nil {
		return "", "", "", "", err
	}
	ioutil.WriteFile(malformed.Name(), malformedConfig, 0644)

	nonValidAuth, err := ioutil.TempFile(dir, "okConfig.yaml")
	if err != nil {
		return "", "", "", "", err
	}
	ioutil.WriteFile(nonValidAuth.Name(), nonValidAuthConfig, 0644)

	return dir, valid.Name(), malformed.Name(), nonValidAuth.Name(), nil
}

func TestLoadConfig(t *testing.T) {
	dir, valid, malformed, nonValidAuth, err := createTempConfigs()
	if err != nil {
		t.Errorf("Config.LoadConfig() Error while creating temp configs %v", err)
	}
	defer os.RemoveAll(dir)

	type args struct {
		configLocation string
		nilLocation    bool
	}

	tests := []struct {
		name string
		args
		compare bool
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				configLocation: valid,
			},
			compare: true,
			wantErr: false,
		},
		{
			name: "non-existent path",
			args: args{
				configLocation: path.Join("config.yaml"),
			},
			wantErr: true,
		},
		{
			name: "no config file provided",
			args: args{
				configLocation: "",
				nilLocation:    true,
			},
			wantErr: true,
		},
		{
			name: "malformed config",
			args: args{
				configLocation: malformed,
			},
			wantErr: true,
		},
		{
			name: "non valid auth config",
			args: args{
				configLocation: nonValidAuth,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg *Config
			var err error

			if tt.args.nilLocation {
				cfg, err = LoadConfig()
			} else {
				cfg, err = LoadConfig(tt.args.configLocation)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.compare {
				equal := reflect.DeepEqual(cfg, exampleConfig)
				if equal == false {
					t.Errorf("Config.LoadConfig() equal = %v, want = true", equal)
					return
				}
			}
		})
	}
}

func TestConfig_Merge(t *testing.T) {
	cfg := Config{
		Endpoint: "https://floodgate",
	}
	cfg2 := Config{
		Endpoint: "https://floodgate2",
	}
	cfg.Merge(cfg2)

	if cfg.Endpoint != cfg2.Endpoint {
		t.Errorf("Config.Merge() endpoint = %s, want %s", cfg.Endpoint, cfg2.Endpoint)
	}
}

func TestSaveConfig(t *testing.T) {
	dir, valid, _, _, err := createTempConfigs()
	if err != nil {
		t.Errorf("Config.LoadConfig() Error while creating temp configs %v", err)
	}
	defer os.RemoveAll(dir)

	cfg, _ := LoadConfig(valid)

	err = SaveConfig(cfg)
	if err != nil {
		t.Errorf("Config.SaveConfig() error = %v, wantErr %v", err, false)
	}
}

var validConfig = []byte("endpoint: https://127.0.0.1/api/v1\ninsecure: true\nauth:\n  basic:\n    enabled: true\n    user: admin\n    password: VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi\nlibraries:\n  - sponnet\nresources:\n  - resources/applications\n  - resources/pipelines\n  - resources/pipelinetemplates")
var malformedConfig = []byte("endpoint: https://127.0.0.1/api/v1\ninsecure: true\nauth:\n  basic:\n    enabled: true\n    user\n    password: VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi\nlibraries:\n  - sponnet\nresources:\n  - resources/applications\n  - resources/pipelines\n  - resources/pipelinetemplates")
var nonValidAuthConfig = []byte("endpoint: https://127.0.0.1/api/v1\ninsecure: true\nauth:\n  basic:\n    enabled: true\n    user: admin\n    password: VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi\n  x509:\n    certPath: ~/.config/floodgate/floodgate.crt\n    enabled: true\n    keyPath: ~/.config/floodgate/floodgate.key\nlibraries:\n  - sponnet\nresources:\n  - resources/applications\n  - resources/pipelines\n  - resources/pipelinetemplates")
