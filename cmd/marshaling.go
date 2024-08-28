package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// jsonMarshalUnmarshal marshals the src value to JSON and then unmarshals it into the dest value.
// dest must be a pointer.
func jsonMarshalUnmarshal(dest any, src any) error {
	// Validate dest is a pointer
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dest)
}
