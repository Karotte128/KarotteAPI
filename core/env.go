package core

import (
	"os"
	"regexp"
)

// regex pattern to replace ${VAR} or ${VAR:-default} in the config
var envPattern = regexp.MustCompile(`\$\{([^}:]+)(:-([^}]+))?\}`)

// expandEnv replaces ${VAR} or ${VAR:-default} with the ENV var in the provided string
func expandEnv(input string) string {
	// search for all matches in input, then execute the replace func on it
	return envPattern.ReplaceAllStringFunc(input, func(match string) string {
		// find all parts
		parts := envPattern.FindStringSubmatch(match)

		// key is the [1] value of the parts
		key := parts[1]
		// defaultValue is the [3] value of the parts
		defaultValue := parts[3]

		// attempt to get the ENV var from key
		val, ok := os.LookupEnv(key)
		if ok {
			// return ENV var value if found
			return val
		}

		// return defaultValue if ENV var is not set
		return defaultValue
	})
}

// expandEnv replaces the ENV vars recursively.
func expandEnvRecursive(input any) any {
	// handle the different possible input types
	switch val := input.(type) {

	// if the type is string, return the result of expandEnv()
	case string:
		return expandEnv(val)

	// if the type is map, run expandEnvRecursive on every value of the map and then return the expanded map
	case map[string]any:
		for k, v := range val {
			val[k] = expandEnvRecursive(v)
		}
		return val

	// if the type is array, run expandEnvRecursive on every element of the array and return the expanded array
	case []any:
		for i, v := range val {
			val[i] = expandEnvRecursive(v)
		}
		return val

	// if the type is none of the above, return the input without expansion
	default:
		return input
	}
}
