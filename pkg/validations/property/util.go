package property

import (
	"encoding/json"
	"fmt"
)

func ConfigToType[T any](in map[string]interface{}, out *T) error {
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshalling input to JSON: %w", err)
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return fmt.Errorf("unmarshalling input to output: %w", err)
	}

	return nil
}
