package fileloader

import (
	"reflect"
	"testing"

	"github.com/google/go-jsonnet"
)

func TestJsonnetLoader_LoadFile(t *testing.T) {
	type fields struct {
		VM *jsonnet.VM
	}
	type args struct {
		filePath string
	}
	vmWithLib := jsonnet.MakeVM()
	vmWithLib.Importer(&jsonnet.FileImporter{JPaths: []string{"testdata/testjsonnetlib.libsonnet"}})
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "load .jsonnet file",
			fields: fields{
				VM: vmWithLib,
			},
			args: args{
				filePath: "testdata/testfile.jsonnet",
			},
			want:    testJsonnetFile,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jl := &JsonnetLoader{
				VM: tt.fields.VM,
			}
			got, err := jl.LoadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonnetLoader.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonnetLoader.LoadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonnetLoader_SupportedFileExtensions(t *testing.T) {
	type fields struct {
		VM *jsonnet.VM
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "support .jsonnet file extension",
			want: []string{".jsonnet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jl := &JsonnetLoader{
				VM: tt.fields.VM,
			}
			if got := jl.SupportedFileExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonnetLoader.SupportedFileExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testJsonnetFile = map[string]interface{}{
	"name":     "testjsonnet",
	"variable": false,
}
