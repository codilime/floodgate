package util

import (
	"testing"
)

func TestAssertMapKeyExists(t *testing.T) {
	tests := []struct {
		name    string
		testMap map[string]interface{}
		key     string
		wantErr bool
	}{
		{
			name:    "key does not exist",
			testMap: map[string]interface{}{},
			key:     "test",
			wantErr: true,
		},
		{
			name: "key exists",
			testMap: map[string]interface{}{
				"test": "test",
			},
			key:     "test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssertMapKeyExists(tt.testMap, tt.key)
			if (got != nil) != tt.wantErr {
				t.Errorf("got %q, wantErr: %v", got, tt.wantErr)
			}
		})
	}
}

func TestAssertMapKeyIsString(t *testing.T) {
	tests := []struct {
		name     string
		testMap  map[string]interface{}
		key      string
		notEmpty bool
		wantErr  bool
	}{
		{
			name:     "key does not exist",
			testMap:  map[string]interface{}{},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists but it's not a string",
			testMap: map[string]interface{}{
				"test": 0,
			},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists and it's a string",
			testMap: map[string]interface{}{
				"test": "",
			},
			key:      "test",
			notEmpty: false,
			wantErr:  false,
		},
		{
			name: "key exists, but it's an empty string",
			testMap: map[string]interface{}{
				"test": "",
			},
			key:      "test",
			notEmpty: true,
			wantErr:  true,
		},
		{
			name: "key exists and it's a non-empty string",
			testMap: map[string]interface{}{
				"test": "test",
			},
			key:      "test",
			notEmpty: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssertMapKeyIsString(tt.testMap, tt.key, tt.notEmpty)
			if (got != nil) != tt.wantErr {
				t.Errorf("got %q, wantErr: %v", got, tt.wantErr)
			}
		})
	}
}

func TestAssertMapKeyIsStringMap(t *testing.T) {
	tests := []struct {
		name     string
		testMap  map[string]interface{}
		key      string
		notEmpty bool
		wantErr  bool
	}{
		{
			name:     "key does not exist",
			testMap:  map[string]interface{}{},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists but it's not a string map",
			testMap: map[string]interface{}{
				"test": 0,
			},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists and it's a string map",
			testMap: map[string]interface{}{
				"test": map[string]interface{}{},
			},
			key:      "test",
			notEmpty: false,
			wantErr:  false,
		},
		{
			name: "key exists, but it's an empty map",
			testMap: map[string]interface{}{
				"test": map[string]interface{}{},
			},
			key:      "test",
			notEmpty: true,
			wantErr:  true,
		},
		{
			name: "key exists and it's not an empty map",
			testMap: map[string]interface{}{
				"test": map[string]interface{}{
					"test": "test",
				},
			},
			key:      "test",
			notEmpty: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssertMapKeyIsStringMap(tt.testMap, tt.key, tt.notEmpty)
			if (got != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", got, tt.wantErr)
			}
		})
	}
}

func TestAssertMapKeyIsInterfaceArray(t *testing.T) {
	tests := []struct {
		name     string
		testMap  map[string]interface{}
		key      string
		notEmpty bool
		wantErr  bool
	}{
		{
			name:     "key does not exist",
			testMap:  map[string]interface{}{},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists but it's not an interface array",
			testMap: map[string]interface{}{
				"test": 0,
			},
			key:      "test",
			notEmpty: false,
			wantErr:  true,
		},
		{
			name: "key exists and it's an interface array",
			testMap: map[string]interface{}{
				"test": []interface{}{},
			},
			key:      "test",
			notEmpty: false,
			wantErr:  false,
		},
		{
			name: "key exists, but it's an empty interface array",
			testMap: map[string]interface{}{
				"test": []interface{}{},
			},
			key:      "test",
			notEmpty: true,
			wantErr:  true,
		},
		{
			name: "key exists and it's not an empty interface array",
			testMap: map[string]interface{}{
				"test": []interface{}{
					nil,
				},
			},
			key:      "test",
			notEmpty: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssertMapKeyIsInterfaceArray(tt.testMap, tt.key, tt.notEmpty)
			if (got != nil) != tt.wantErr {
				t.Errorf("got error %q, wantErr: %v", got, tt.wantErr)
			}
		})
	}
}
