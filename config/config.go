package config

import (
	"encoding/json"
	"io"
)

// Config type
type Config map[string]interface{}

// LoadConfig will load config into Config type
func LoadConfig(s io.Reader) (Config, error) {
	var config Config

	decoder := json.NewDecoder(s)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
