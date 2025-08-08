// Package file implements filesystem handling functions.
package file

import (
	"encoding/json"
	"os"
)

// ReadJSON unmarshals a JSON file into an object.
func ReadJSON(orig string, data any) error {
	bytes, err := os.ReadFile(orig)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}
