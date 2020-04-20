package fileloader

import (
	"encoding/json"
	"io/ioutil"

	"github.com/google/go-jsonnet"
)

// NewJsonnetLoader create loader
func NewJsonnetLoader(librariesPaths ...string) *JsonnetLoader {
	parser := &JsonnetLoader{}
	parser.VM = jsonnet.MakeVM()
	parser.VM.Importer(&jsonnet.FileImporter{JPaths: librariesPaths})
	return parser
}

// JsonnetLoader load jsonnet files
type JsonnetLoader struct {
	*jsonnet.VM
}

// LoadFile load file
func (jl *JsonnetLoader) LoadFile(filePath string) ([]map[string]interface{}, error) {
	inputFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	evaluatedSnippet, err := jl.EvaluateSnippetStream(filePath, string(inputFile))
	if err != nil {
    evaluatedSingleSnippet, err := jl.EvaluateSnippet(filePath, string(inputFile))
    if err != nil {
      return nil, err
    }
    evaluatedSnippet = []string{evaluatedSingleSnippet}
	}
	var output []map[string]interface{}
  for i := range evaluatedSnippet {
    var partial map[string]interface{}
    json.Unmarshal([]byte(evaluatedSnippet[i]), &partial)
    output = append(output, partial)
  }
	return output, nil
}

// SupportedFileExtensions get list of supported file extensions
func (jl *JsonnetLoader) SupportedFileExtensions() []string {
	return []string{".jsonnet"}
}
