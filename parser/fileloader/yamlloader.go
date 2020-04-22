package fileloader

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// NewYAMLLoader create new YAML file loader
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// YAMLLoader load YAML files
type YAMLLoader struct{}

// LoadFile load file
func (yl *YAMLLoader) LoadFile(filePath string) ([]map[string]interface{}, error) {
	inputFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	if err := yaml.Unmarshal(inputFile, &output); err != nil {
		return nil, err
	}
	return []map[string]interface{}{output}, nil
}

// SupportedFileExtensions get supported file extensions
func (yl *YAMLLoader) SupportedFileExtensions() []string {
	return []string{".yaml", ".yml"}
}
