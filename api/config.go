package api

import (
	cfg "github.com/karotte128/karottelib/config"
)

// Config contains all details to create a new api.
type Config map[string]any

// Global config instance
var config Config

// loadConfig sets the global config instance.
func loadConfig(conf Config) {
	config = conf
}

// GetModuleConfig returns the raw config block for a module.
func GetModuleConfig(moduleName string) (Config, bool) {
	if config == nil {
		return nil, false
	}

	return cfg.GetNestedValue[map[string]any](config, "modules", moduleName)
}

// GetMiddlewareConfig returns the raw config block for a module.
func GetMiddlewareConfig(middlewareName string) (Config, bool) {
	if config == nil {
		return nil, false
	}

	return cfg.GetNestedValue[map[string]any](config, "middleware", middlewareName)
}

// GetServerConfig returns the server config.
func getServerConfig() (Config, bool) {
	if config == nil {
		return nil, false
	}

	return cfg.GetNestedValue[map[string]any](config, "server")
}
