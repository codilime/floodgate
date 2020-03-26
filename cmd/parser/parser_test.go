package parser

import (
	"reflect"
	"testing"
)

func TestParser_loadFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name          string
		librariesPath []string
		args          args
		want          map[string]interface{}
		wantErr       bool
	}{
		{
			name:          "load jsonnet file",
			librariesPath: []string{"testdata/testsimplelib"},
			args: args{
				filePath: "testdata/testsimple/testsimple.jsonnet",
			},
			want: map[string]interface{}{
				"name":     "testjsonnet",
				"variable": false,
			},
			wantErr: false,
		},
		{
			name:          "load json file",
			librariesPath: []string{},
			args: args{
				filePath: "testdata/testsimple/testsimple.json",
			},
			want: map[string]interface{}{
				"name":     "testjson",
				"variable": false,
			},
			wantErr: false,
		},
		{
			name:          "load yaml file",
			librariesPath: []string{},
			args: args{
				filePath: "testdata/testsimple/testsimple.yaml",
			},
			want: map[string]interface{}{
				"name":     "testyaml",
				"variable": false,
			},
			wantErr: false,
		},
		{
			name:          "load yml file",
			librariesPath: []string{},
			args: args{
				filePath: "testdata/testsimple/testsimple.yml",
			},
			want: map[string]interface{}{
				"name":     "testyml",
				"variable": false,
			},
			wantErr: false,
		},
		{
			name:          "fail to load file with invalid extension",
			librariesPath: []string{},
			args: args{
				filePath: "testdata/testsimple/testsimple.invalid_ext",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:          "fail to load nonexistent file",
			librariesPath: []string{},
			args: args{
				filePath: "testdata/testsimple/thefilethatshouldnotbe",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateParser(tt.librariesPath)
			got, err := p.loadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.loadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.loadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			librariesPath: []string{"testdata/testsimplelib"},
			args:          args{entrypoint: "testdata/testsimple"},
			want:          testSimpleResources,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateParser(tt.librariesPath)
			got, err := p.loadDirectory(tt.args.entrypoint)
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

func TestParser_LoadObjectsFromDirectories(t *testing.T) {
	type args struct {
		directories []string
	}
	tests := []struct {
		name          string
		librariesPath []string
		args          args
		wantErr       bool
	}{
		{
			name:          "successfully load objects from single directory",
			librariesPath: []string{},
			args: args{
				directories: []string{"testdata/testapplication"},
			},
			wantErr: false,
		},
		{
			name:          "successfully load objects from multiple directories",
			librariesPath: []string{},
			args: args{
				directories: []string{"testdata/testapplication", "testdata/testpipeline"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateParser(tt.librariesPath)
			if err := p.LoadObjectsFromDirectories(tt.args.directories); (err != nil) != tt.wantErr {
				t.Errorf("Parser.LoadObjectsFromDirectories() error = %v, wantErr %v", err, tt.wantErr)
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
				objects: []map[string]interface{}{testApplication},
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
			if err := p.readObjects(tt.args.objects); (err != nil) != tt.wantErr {
				t.Errorf("Parser.readObjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testSimpleResources = []map[string]interface{}{
	map[string]interface{}{
		"name":     "testjson",
		"variable": false,
	},
	map[string]interface{}{
		"name":     "testjsonnet",
		"variable": false,
	},
	map[string]interface{}{
		"name":     "testyaml",
		"variable": false,
	},
	map[string]interface{}{
		"name":     "testyml",
		"variable": false,
	},
}

var testApplication = map[string]interface{}{
	"name":           "testapplication",
	"description":    "Test application",
	"email":          "test@floodgate.com",
	"user":           "test@floodgate.com",
	"cloudProviders": "kubernetes",
	"dataSources": map[string]interface{}{
		"disabled": []string{},
		"enabled":  []string{},
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
	"trafficGuards": []string{},
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
