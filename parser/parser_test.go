package parser

import (
	"reflect"
	"testing"
)

func TestParser_loadDirectory(t *testing.T) {
	type args struct {
		entrypoint string
	}
	tests := []struct {
		name          string
		librariesPath []string
		args          args
		want          []map[string]interface{}
		wantErr       bool
	}{
		{
			name:          "succesfully load directory",
			librariesPath: []string{"testdata/testlibraries"},
			args:          args{entrypoint: "testdata/applications"},
			want:          []map[string]interface{}{testApplicationJSON, testApplicationJsonnet},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewResourceParser(tt.librariesPath...)
			got, err := p.loadFilesFromDirectory(tt.args.entrypoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.loadDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.loadDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ParseDirectories(t *testing.T) {
	type args struct {
		directories []string
	}
	tests := []struct {
		name          string
		librariesPath []string
		args          args
		want          *ParsedResourceData
		wantErr       bool
	}{
		{
			name:          "successfully load objects from single directory",
			librariesPath: []string{"testdata/testlibraries"},
			args: args{
				directories: []string{"testdata/applications"},
			},
			want:    &ParsedResourceData{Applications: []map[string]interface{}{testApplicationJSON, testApplicationJsonnet}},
			wantErr: false,
		},
		{
			name:          "successfully load objects from multiple directories",
			librariesPath: []string{"testdata/testlibraries"},
			args: args{
				directories: []string{"testdata/applications", "testdata/pipelines"},
			},
			want:    &ParsedResourceData{Applications: []map[string]interface{}{testApplicationJSON, testApplicationJsonnet}, Pipelines: []map[string]interface{}{testPipelineJSON}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewResourceParser(tt.librariesPath...)
			if err != nil {
				t.Error(err)
			}
			got, err := p.ParseDirectories(tt.args.directories)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseDirectories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.loadDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_readObjects(t *testing.T) {
	type args struct {
		objects []map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "can't read empty list of objects",
			args: args{
				objects: []map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "read application",
			args: args{
				objects: []map[string]interface{}{testApplicationJSON},
			},
			wantErr: false,
		},
		{
			name: "read pipeline template",
			args: args{
				objects: []map[string]interface{}{testPipelineTemplate},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			if _, err := p.parseObjects(tt.args.objects); (err != nil) != tt.wantErr {
				t.Errorf("Parser.readObjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testApplicationJSON = map[string]interface{}{
	"name":           "testappjson",
	"description":    "Test application from JSON file.",
	"email":          "example@example.com",
	"user":           "example@example.com",
	"cloudProviders": "kubernetes",
	"dataSources": map[string]interface{}{
		"disabled": []interface{}{},
		"enabled":  []interface{}{},
	},
	"platformHealthOnly":             false,
	"platformHealthOnlyShowOverride": false,
	"providerSettings": map[string]interface{}{
		"aws": map[string]interface{}{
			"useAmiBlockDeviceMappings": false,
		},
		"gce": map[string]interface{}{
			"associatePublicIpAddress": false,
		},
	},
	"trafficGuards": []interface{}{},
}

var testApplicationJsonnet = map[string]interface{}{
	"cloudProviders": "kubernetes",
	"dataSources": map[string]interface{}{
		"disabled": []interface{}{},
		"enabled":  []interface{}{},
	},
	"description":                    "Test application from Jsonnet file.",
	"email":                          "example@example.com",
	"name":                           "testappjsonnet",
	"platformHealthOnly":             false,
	"platformHealthOnlyShowOverride": false,
	"providerSettings": map[string]interface{}{
		"aws": map[string]interface{}{
			"useAmiBlockDeviceMappings": false,
		},
		"gce": map[string]interface{}{
			"associatePublicIpAddress": false,
		},
	},
	"trafficGuards": []interface{}{},
	"user":          "example@example.com",
}

var testPipelineJSON = map[string]interface{}{
	"application":          "testpipelineapplication",
	"keepWaitingPipelines": false,
	"limitConcurrent":      true,
	"name":                 "testpipelinejson",
	"notifications":        []interface{}{},
	"stages":               []interface{}{},
	"triggers":             []interface{}{},
}

var testPipelineTemplate = map[string]interface{}{
	"id": "testpipelinetemplate",
	"metadata": map[string]interface{}{
		"name":        "Test pipeline template",
		"description": "Test pipeline template.",
		"owner":       "example@example.com",
		"scopes": []interface{}{
			"",
		},
	},
	"variables": []map[string]interface{}{},
	"schema":    "v2",
}
