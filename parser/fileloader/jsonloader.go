package fileloader

import (
	"encoding/json"
	"io/ioutil"
)

// NewJSONLoader create loader that can load JSON files with extension ".json"
func NewJSONLoader() *JSONLoader {
	return &JSONLoader{}
}

// JSONLoader load JSON files
type JSONLoader struct{}

// LoadFile load file
func (jl *JSONLoader) LoadFile(filePath string) (map[string]interface{}, error) {
	inputFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	if err := json.Unmarshal(inputFile, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// SupportedFileExtensions get supported file extensions
func (jl *JSONLoader) SupportedFileExtensions() []string {
	return []string{".json"}
}
