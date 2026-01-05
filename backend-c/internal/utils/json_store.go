package utils

import (
	"encoding/json"
	"os"
)

// Read JSON from file
func ReadJSON(file string, data interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, data)
}

// Write JSON to file
func WriteJSON(file string, data interface{}) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, content, 0644)
}
