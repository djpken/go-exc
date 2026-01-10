// Package utils provides utility functions for OKEx exchange
package utils

import (
	"encoding/json"
)

// S2M converts a struct to map[string]string
func S2M(i interface{}) map[string]string {
	m := make(map[string]string)
	j, _ := json.Marshal(i)
	_ = json.Unmarshal(j, &m)

	return m
}
