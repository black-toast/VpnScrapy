package util

import (
	"encoding/json"
)

func Parse[T any](jsonData []byte, instance T) (T, error) {
	if err := json.Unmarshal(jsonData, &instance); err == nil {
		return instance, err
	}
	return instance, nil
}
