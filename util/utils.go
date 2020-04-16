package util

import (
	"fmt"
	"os"
)

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

// CreateDirs create directories
func CreateDirs(dirPaths ...string) error {
	for _, dirPath := range dirPaths {
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
