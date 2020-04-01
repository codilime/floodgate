package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
)

// ResourceData contains all the managed ResourceData defined for Floodgate.
type ResourceData struct {
	Applications      []map[string]interface{}
	Pipelines         []map[string]interface{}
	PipelineTemplates []map[string]interface{}
}

// Parser extends jsonnet VM with floodgate-specific configuration.
type Parser struct {
	*jsonnet.VM
	Resources ResourceData
}

// CreateParser creates an instance of Floodgate resources parser.
func CreateParser(librariesPath []string) *Parser {
	parser := &Parser{}
	parser.VM = jsonnet.MakeVM()
	parser.VM.Importer(&jsonnet.FileImporter{
		JPaths: librariesPath,
	})
	return parser
}

func (p *Parser) loadFile(filePath string) (map[string]interface{}, error) {
	fileExt := filepath.Ext(filePath)
	switch fileExt {
	case ".jsonnet":
		return p.loadJsonnetFile(filePath)
	case ".json":
		return p.loadJSONFile(filePath)
	case ".yaml", ".yml":
		return p.loadYAMLFile(filePath)
	default:
		return nil, fmt.Errorf("unsupported file extension %q", fileExt)
	}
}

func (p *Parser) loadJsonnetFile(filePath string) (map[string]interface{}, error) {
	inputFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	parserOutput, err := p.EvaluateSnippet(filePath, string(inputFile))
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	json.Unmarshal([]byte(parserOutput), &output)
	return output, nil
}

func (p *Parser) loadJSONFile(filePath string) (map[string]interface{}, error) {
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

func (p *Parser) loadYAMLFile(filePath string) (map[string]interface{}, error) {
	inputFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	if err := yaml.Unmarshal(inputFile, &output); err != nil {
		return nil, err
	}
	return output, nil
}

func (p *Parser) loadDirectory(entrypoint string) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	err := filepath.Walk(entrypoint,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			obj, err := p.loadFile(path)
			// We want to load only supported files, but don't want to stop loading after
			// finding unsupported one
			if err != nil {
				log.Printf("%s not loaded due to %v", f.Name(), err)
				return nil
			}
			objects = append(objects, obj)
			log.Printf("%s loaded", f.Name())
			return nil
		})
	if err != nil {
		return nil, err
	}
	return objects, nil
}

// LoadObjectsFromDirectories walks through provided directories to parse provided configuration files
// and catalogs them according to their types.
func (p *Parser) LoadObjectsFromDirectories(directories []string) error {
	var objects []map[string]interface{}
	for _, entrypoint := range directories {
		output, err := p.loadDirectory(entrypoint)
		if err != nil {
			log.Fatal(err)
			return err
		}
		objects = append(objects, output...)
	}
	err := p.readObjects(objects)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) readObjects(objects []map[string]interface{}) error {
	// TODO(wurbanski): Verify heuristics for determining spinnaker object types
	if len(objects) == 0 {
		return fmt.Errorf("no objects found")
	}
	for _, object := range objects {
		if _, ok := object["application"]; ok {
			p.Resources.Pipelines = append(p.Resources.Pipelines, object)
			continue
		}
		if _, ok := object["variables"]; ok {
			p.Resources.PipelineTemplates = append(p.Resources.PipelineTemplates, object)
			continue
		}
		if _, ok := object["providerSettings"]; ok {

			p.Resources.Applications = append(p.Resources.Applications, object)
			continue
		}
		return fmt.Errorf("object %v not of any known type", object)
	}
	return nil
}
