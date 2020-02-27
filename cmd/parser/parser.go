package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

// SpinnakerResources contains all the managed SpinnakerResources defined for Floodgate.
type SpinnakerResources struct {
	Applications      []map[string]interface{}
	Pipelines         []map[string]interface{}
	PipelineTemplates []map[string]interface{}
}

// Parser extends jsonnet VM with floodgate-specific configuration.
type Parser struct {
	*jsonnet.VM
	Resources SpinnakerResources
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

func (p *Parser) loadDirectory(entrypoint string) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	err := filepath.Walk(entrypoint,
		func(path string, f os.FileInfo, err error) error {
			log.Print(f.Name())
			if err != nil {
				log.Fatal(err)
				return err
			}
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".jsonnet") {
				return nil
			}
			obj, err := p.loadJsonnetFile(path)
			if err != nil {
				log.Fatal(err)
				return err
			}
			objects = append(objects, obj)
			return nil
		})
	if err != nil {
		log.Fatal(err)
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
			err := p.validatePipelineTemplate(object)
			if err != nil {
				return fmt.Errorf("Encountered an error while reading pipeline template %v: %w", object, err)
			}
			p.Resources.PipelineTemplates = append(p.Resources.PipelineTemplates, object)
			continue
		}
		if _, ok := object["cloudProviders"]; ok {

			p.Resources.Applications = append(p.Resources.Applications, object)
			continue
		}
		return fmt.Errorf("object %v not of any known type", object)
	}
	return nil
}

func (p Parser) validatePipelineTemplate(object map[string]interface{}) error {
	if _, ok := object["id"]; ok != true {
		return fmt.Errorf("missing field 'id'")
	}
	metadata, ok := object["metadata"]
	if ok != true {
		return fmt.Errorf("missing field 'metadata'")
	}
	if _, ok := metadata.(map[string]interface{})["name"]; ok != true {
		return fmt.Errorf("missing key 'name' in map 'metadata`")
	}
	if _, ok := metadata.(map[string]interface{})["owner"].(string); ok != true {
		return fmt.Errorf("missing key 'owner' in map 'metadata'")
	}
	if _, ok := metadata.(map[string]interface{})["description"].(string); ok != true {
		return fmt.Errorf("missing key 'description' in map 'metadata'")
	}
	if _, ok := metadata.(map[string]interface{})["scopes"].([]interface{}); ok != true {
		return fmt.Errorf("missing key 'scopes' in map 'metadata'")
	}
	return nil
}
