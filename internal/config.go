package internal

import "github.com/karotte128/karotteapi"

// Global config instance
var config karotteapi.Config

func LoadConfig(conf karotteapi.Config) {
	config = conf
}

// GetModuleConfig returns the raw config block for a module.
func GetModuleConfig(moduleName string) (karotteapi.Config, bool) {
	if config == nil {
		return nil, false
	}

	return GetNestedValue[karotteapi.Config](config, "modules", moduleName)
}

// GetMiddlewareConfig returns the raw config block for a module.
func GetMiddlewareConfig(middlewareName string) (karotteapi.Config, bool) {
	if config == nil {
		return nil, false
	}

	return GetNestedValue[karotteapi.Config](config, "middleware", middlewareName)
}

// GetServerConfig returns the server config.
func GetServerConfig() (karotteapi.Config, bool) {
	if config == nil {
		return nil, false
	}

	return GetNestedValue[karotteapi.Config](config, "server")
}
