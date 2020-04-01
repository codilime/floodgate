package util

import "fmt"

// CombineErrors combine array of errors into one
func CombineErrors(errors []error) error {
	var combinedErr error
	for idx, err := range errors {
		if err != nil {
			combinedErr = err
			errors = errors[idx+1:]
			break
		}
	}
	for _, err := range errors {
		if err != nil {
			combinedErr = fmt.Errorf("%v, %w", combinedErr.Error(), err)
		}
	}
	return combinedErr
}
