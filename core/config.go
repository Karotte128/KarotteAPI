package core

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Global config instance
var config map[string]any

// LoadConfig loads, expands (replace ENV vars), and parses the configuration file.
func LoadConfig(path string) error {
	// Reading the provided file path
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the file contents into a map[string]any
	var raw map[string]any
	err = toml.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// Expand the ENV vars and set global config
	config = expandEnvRecursive(raw).(map[string]any)

	// no error, yay!
	return nil
}

// GetModuleConfig returns the raw config block for a module.
func GetModuleConfig(moduleName string) (map[string]any, bool) {
	if config == nil {
		return nil, false
	}

	return GetNestedValue[map[string]any](config, "modules", moduleName)
}

// GetServerConfig returns the server config.
func GetServerConfig() (map[string]any, bool) {
	if config == nil {
		return nil, false
	}

	return GetNestedValue[map[string]any](config, "server")
}
