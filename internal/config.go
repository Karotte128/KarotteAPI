package internal

// Global config instance
var config map[string]any

func LoadConfig(conf map[string]any) {
	config = conf
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
