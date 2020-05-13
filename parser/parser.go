package parser

import (
	"fmt"
	"os"
	"path/filepath"

	fl "github.com/codilime/floodgate/parser/fileloader"
	log "github.com/sirupsen/logrus"
)

// ParsedResourceData contains all the managed ResourceData defined for Floodgate.
type ParsedResourceData struct {
	Applications      []map[string]interface{}
	Pipelines         []map[string]interface{}
	PipelineTemplates []map[string]interface{}
}

// NewJsonnetParser create parser to parse Jsonnet files
func NewJsonnetParser(librariesPath ...string) (*Parser, error) {
	jsonnetLoader := fl.NewJsonnetLoader(librariesPath...)
	return NewParser(jsonnetLoader)
}

// NewResourceParser create parser to parser JSON, Jsonnet and YAML files
func NewResourceParser(librariesPath ...string) (*Parser, error) {
	jsonLoader := fl.NewJSONLoader()
	jsonnetLoader := fl.NewJsonnetLoader(librariesPath...)
	yamlLoader := fl.NewYAMLLoader()
	return NewParser(jsonLoader, jsonnetLoader, yamlLoader)
}

// NewParser create new parser
func NewParser(fileLoaders ...fl.FileLoader) (*Parser, error) {
	parser := &Parser{fileLoaderByExtension: make(map[string]fl.FileLoader)}
	for _, fileLoader := range fileLoaders {
		if err := parser.RegisterFileLoader(fileLoader); err != nil {
			return nil, err
		}
	}
	return parser, nil
}

// Parser parse files to get resources
type Parser struct {
	fileLoaderByExtension map[string]fl.FileLoader
}

// RegisterFileLoader register new file loaders
func (p *Parser) RegisterFileLoader(fileLoader fl.FileLoader) error {
	for _, fileExt := range fileLoader.SupportedFileExtensions() {
		_, exists := p.fileLoaderByExtension[fileExt]
		if exists {
			return fmt.Errorf("loader for file extension %q was already registered ", fileExt)
		}
		p.fileLoaderByExtension[fileExt] = fileLoader
	}
	return nil
}

// ParseDirectories walks through provided directories to parse provided configuration files
// and catalogs them according to their types.
func (p *Parser) ParseDirectories(directories []string) (*ParsedResourceData, error) {
	var objects []map[string]interface{}
	for _, entrypoint := range directories {
		output, err := p.loadFilesFromDirectory(entrypoint)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		objects = append(objects, output...)
	}
	return p.parseObjects(objects)
}

func (p *Parser) loadFilesFromDirectory(entrypoint string) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	err := filepath.Walk(entrypoint,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			fileExt := filepath.Ext(path)
			fileLoader, ok := p.fileLoaderByExtension[fileExt]
			// We want to load only supported files, but don't want to stop loading after
			// finding unsupported one
			if !ok {
				log.Warn(" unsupported file extension\n", fileExt)
				return nil
			}
			obj, err := fileLoader.LoadFile(path)
			if err != nil {
				log.Warn(f.Name(), " not loaded due to\n", err)
				return nil
			}
			objects = append(objects, obj...)
			log.Debugf("Loaded file: %s", path)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (p *Parser) parseObjects(objects []map[string]interface{}) (*ParsedResourceData, error) {
	parsedResources := &ParsedResourceData{}
	// TODO(wurbanski): Verify heuristics for determining spinnaker object types
	if len(objects) == 0 {
		return nil, fmt.Errorf("no objects found")
	}
	for _, object := range objects {
		if _, ok := object["application"]; ok {
			parsedResources.Pipelines = append(parsedResources.Pipelines, object)
			continue
		}
		if _, ok := object["variables"]; ok {
			parsedResources.PipelineTemplates = append(parsedResources.PipelineTemplates, object)
			continue
		}
		if _, ok := object["email"]; ok {
			parsedResources.Applications = append(parsedResources.Applications, object)
			continue
		}
	}
	return parsedResources, nil
}
