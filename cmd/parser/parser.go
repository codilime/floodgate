package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

type Parser struct {
	*jsonnet.VM
}

func CreateParser(librariesPath []string) (*Parser, error) {
	parser := &Parser{}
	parser.VM = jsonnet.MakeVM()
	parser.VM.Importer(&jsonnet.FileImporter{
		JPaths: librariesPath,
	})
	return parser, nil
}

func (p *Parser) LoadJsonnetFile(filePath string) (map[string]interface{}, error) {
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

func (p *Parser) LoadDirectory(entrypoint string) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	err := filepath.Walk(entrypoint,
		func(path string, f os.FileInfo, err error) error {
			log.Print(f.Name())
			if err != nil {
				log.Fatal(err)
				return err
			}
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".jsonnet") {
				log.Print("skipped")
				return nil
			}
			log.Print("loaded")
			obj, err := p.LoadJsonnetFile(path)
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
	log.Print("Objects:", objects)
	return objects, nil
}

func (p *Parser) LoadDirectories(directories []string) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	for _, entrypoint := range directories {
		output, err := p.LoadDirectory(entrypoint)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		objects = append(objects, output...)
	}
	return objects, nil
}
