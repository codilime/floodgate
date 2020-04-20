package fileloader

import (
	"reflect"
	"testing"
)

func TestJSONLoader_LoadFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "load json file",
			args: args{
				filePath: "testdata/testfile.json",
			},
			want:    testJSONFIle,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jl := &JSONLoader{}
			got, err := jl.LoadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONLoader.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONLoader.LoadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONLoader_SupportedFileExtensions(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "support .json file extension",
			want: []string{".json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jl := &JSONLoader{}
			if got := jl.SupportedFileExtensions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONLoader.SupportedFileExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testJSONFIle = map[string]interface{}{
	"name":     "testjson",
	"variable": false,
}
