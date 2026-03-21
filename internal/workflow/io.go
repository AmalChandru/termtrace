package workflow

import (
	"encoding/json"
	"fmt"
	"os"
)

func Marshal(w *Workflow) ([]byte, error) {
	if err := w.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(w, "", "  ")
}

func Unmarshal(data []byte) (*Workflow, error) {
	var w Workflow
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("workflow: invalid JSON: %w", err)
	}
	if err := w.Validate(); err != nil {
		return nil, err
	}
	return &w, nil
}

// TODO: Uses 0644; tighten it.
func WriteToFile(w *Workflow, path string) error {
	data, err := Marshal(w)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func ReadFromFile(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}
