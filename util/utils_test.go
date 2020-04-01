package util

import (
	"fmt"
	"testing"
)

func TestCombineErrors(t *testing.T) {
	tests := []struct {
		name   string
		errors []error
		want   error
	}{
		{
			name:   "no errors",
			errors: []error{},
			want:   nil,
		},
		{
			name:   "one nil error",
			errors: []error{nil},
			want:   nil,
		},
		{
			name:   "multiple nil errors",
			errors: []error{nil, nil},
			want:   nil,
		},
		{
			name:   "one error",
			errors: []error{fmt.Errorf("test error")},
			want:   fmt.Errorf("test error"),
		},
		{
			name:   "error followed by nil error",
			errors: []error{fmt.Errorf("test error"), nil},
			want:   fmt.Errorf("test error"),
		},
		{
			name:   "nil error followed by error",
			errors: []error{nil, fmt.Errorf("test error")},
			want:   fmt.Errorf("test error"),
		},
		{
			name:   "multiple errors",
			errors: []error{fmt.Errorf("first test error"), fmt.Errorf("second test error")},
			want:   fmt.Errorf("first test error, second test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CombineErrors(tt.errors)
			assertError(t, got, tt.want)
		})
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	if want == nil {
		assertNoError(t, got)
		return
	}
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want.Error() {
		t.Errorf("got %q, want %q", got, want)
	}
}
