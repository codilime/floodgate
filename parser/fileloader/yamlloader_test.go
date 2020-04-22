package fileloader

import (
	"reflect"
	"testing"
)

func TestYAMLLoader_LoadFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "load .yaml file",
			args: args{
				filePath: "testdata/testfile.yaml",
			},
			want:    testYAMLFile,
			wantErr: false,
		},
		{
			name: "load .yml file",
			args: args{
				filePath: "testdata/testfile.yml",
			},
			want:    testYMLFile,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yl := &YAMLLoader{}
			got, err := yl.LoadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("YAMLLoader.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YAMLLoader.LoadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYAMLLoader_SupportedFileExtensions(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "support .yaml and .yml file extensions",
			want: []string{".yaml", ".yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yl := &YAMLLoader{}
			if got := yl.SupportedFileExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YAMLLoader.SupportedFileExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testYAMLFile = []map[string]interface{}{{
	"name":     "testyaml",
	"variable": false,
}}

var testYMLFile = []map[string]interface{}{{
	"name":     "testyml",
	"variable": false,
}}
