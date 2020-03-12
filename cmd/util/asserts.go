package util

import (
	"fmt"
	"strings"
)

// AssertMapKeyExists check wheter or not map contains a key
func AssertMapKeyExists(object map[string]interface{}, key string) error {
	_, exists := object[key]
	if exists != true {
		return fmt.Errorf("required key '%v' missing", key)
	}
	return nil
}

// AssertMapKeyIsString check wheter or not map contains a key
func AssertMapKeyIsString(object map[string]interface{}, key string, notEmpty bool) error {
	if err := AssertMapKeyExists(object, key); err != nil {
		return err
	}
	value, ok := object[key].(string)
	if ok != true {
		return fmt.Errorf("required key '%v' must be a string", key)
	}
	if notEmpty && strings.TrimSpace(value) == "" {
		return fmt.Errorf("required key '%v' must not be empty or be a string consisting only of spaces", key)
	}
	return nil
}

// AssertMapKeyIsStringMap check whether or not map contains a key
func AssertMapKeyIsStringMap(object map[string]interface{}, key string, notEmpty bool) error {
	if err := AssertMapKeyExists(object, key); err != nil {
		return err
	}
	value, ok := object[key].(map[string]interface{})
	if ok != true {
		return fmt.Errorf("required key '%v' must be a map", key)
	}
	if notEmpty && len(value) == 0 {
		return fmt.Errorf("required key '%v' must not be empty", key)
	}
	return nil
}

// AssertMapKeyIsInterfaceArray check whether or not map contains a key
func AssertMapKeyIsInterfaceArray(object map[string]interface{}, key string, notEmpty bool) error {
	if err := AssertMapKeyExists(object, key); err != nil {
		return err
	}
	value, ok := object[key].([]interface{})
	if ok != true {
		return fmt.Errorf("required key '%v' must be an array", key)
	}
	if notEmpty && len(value) == 0 {
		return fmt.Errorf("required key '%v' must not be empty", key)
	}
	return nil
}
