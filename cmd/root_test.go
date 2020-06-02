package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

func CreateTempFiles(endpoint string, createSponnet bool) (string, string, error) {
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		return "", "", err
	}

	//Create applications dir with example json app declaration
	applicationsDir, err := ioutil.TempDir(dir, "applications")
	if err != nil {
		return "", "", err
	}

	application, err := ioutil.TempFile(applicationsDir, "*.jsonapp.json")
	if err != nil {
		return "", "", err
	}
	ioutil.WriteFile(application.Name(), []byte("{\n   \"cloudProviders\": \"kubernetes\",\n   \"dataSources\": {\n      \"disabled\": [],\n      \"enabled\": []\n   },\n   \"description\": \"Example application created from JSON file.\",\n   \"email\": \"example@example.com\",\n   \"name\": \"jsonapp\",\n   \"platformHealthOnly\": false,\n   \"platformHealthOnlyShowOverride\": false,\n   \"providerSettings\": {\n      \"aws\": {\n         \"useAmiBlockDeviceMappings\": false\n      },\n      \"gce\": {\n         \"associatePublicIpAddress\": false\n      }\n   },\n   \"trafficGuards\": [],\n   \"user\": \"floodgate@example.com\"\n}"), 0644)

	//Create pipelines dir with example json pipeline declaration
	pipelinesDir, err := ioutil.TempDir(dir, "pipelines")
	if err != nil {
		return "", "", err
	}

	if createSponnet {
		pipelineJsonnet, err := ioutil.TempFile(pipelinesDir, "*.jsonnetpipeline.jsonnet")
		if err != nil {
			return "", "", err
		}
		ioutil.WriteFile(pipelineJsonnet.Name(), []byte("local pipelines = import 'pipeline.libsonnet';\n\npipelines.pipeline()\n.withName('Example pipeline from Jsonnet')\n.withId('jsonnetpipeline')\n.withApplication('jsonnetapp')"), 0644)

	}
	pipeline, err := ioutil.TempFile(pipelinesDir, "*.jsonpipeline.json")
	if err != nil {
		return "", "", err
	}
	ioutil.WriteFile(pipeline.Name(), []byte("{\n   \"application\": \"jsonnetapp\",\n   \"id\": \"jsonnetpipeline\",\n   \"keepWaitingPipelines\": false,\n   \"limitConcurrent\": true,\n   \"name\": \"Example pipeline from Jsonnet\",\n   \"notifications\": [ ],\n   \"stages\": [ ],\n   \"triggers\": [ ]\n}\n"), 0644)

	//Create pipeline templates dir with example json pipeline template declaration
	pipelineTemplatesDir, err := ioutil.TempDir(dir, "pipelinetemplates")
	if err != nil {
		return "", "", err
	}

	pipelineTemplate, err := ioutil.TempFile(pipelineTemplatesDir, "*.jsonpipeline.json")
	if err != nil {
		return "", "", err
	}
	ioutil.WriteFile(pipelineTemplate.Name(), []byte("{\n   \"id\": \"jsonnetpt\",\n   \"metadata\": {\n      \"description\": \"Example pipeline template created from Jsonnet file.\",\n      \"name\": \"Example pipeline template from Jsonnet\",\n      \"owner\": \"floodgate@example.com\",\n      \"scopes\": [\n         \"global\"\n      ]\n   },\n   \"protect\": false,\n   \"schema\": \"v2\",\n   \"variables\": [ ]\n}"), 0644)

	//Get sponnet directory
	_, filename, _, _ := runtime.Caller(0)
	sponnetDir := path.Join(path.Dir(filename), "../sponnet")

	//Create config file
	config, err := ioutil.TempFile(dir, "*.config.yaml")
	if err != nil {
		return "", "", err
	}
	configStr := fmt.Sprintf("endpoint: %s\ninsecure: true\nauth:\n  basic:\n    enabled: true\n    user: admin\n    password: VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi\nlibraries:\n  - %s\nresources:\n  - %s\n  - %s\n  - %s", endpoint, sponnetDir, applicationsDir, pipelinesDir, pipelineTemplatesDir)
	ioutil.WriteFile(config.Name(), []byte(configStr), 0644)

	return dir, config.Name(), nil
}

func RemoveTempDir(dir string) {
	if dir != "" && strings.HasPrefix(dir, "/tmp/testing") {
		os.RemoveAll(dir)
	}
}
