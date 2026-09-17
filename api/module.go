package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"

	cfg "github.com/karotte128/karottelib/config"
)

// A module is a component handles requests to an API endpoint.
// If a module sets the prefix "/example/", all requests to "/example/*" get handled by the module.
//
// A module describes:
// - a unique name (e.g., "status")
// - a set of routes served by an http.NewServeMux()
// - a startup function (optional)
// - a shutdown function (optional)

// Modules register themselves automatically via init() inside their
// own package. The core does not need to know about them explicitly.

// Module is the struct the module needs to provide to the module registry to register itself.
type Module struct {
	// Name is the name of the module. It is used for logging.
	Name string

	// Routes returns a URL prefix and an http.Handler that serves all routes
	// for this module.
	//
	// Example:
	//   prefix = "/example/"
	//   handler = http.HandlerFunc()
	Routes func() (prefix string, handler http.Handler)

	// Startup is a function that is run on startup.
	// This can be used to initialize a connection to external services like databases.
	Startup func() error

	// Shutdown is a function that is run on shutdown.
	// This can be used to cleanly disconnect from services connected during startup.
	Shutdown func() error
}

// ModuleStatus holds status information about all registered modules.
// It contains lists of modules, sorted by status type.
type ModuleStatus struct {
	ModuleCount     int
	RunningModules  []string
	DisabledModules []string
	FailedModules   []string
}

// moduleRegistry holds all globally registered modules.
var moduleRegistry []Module

// moduleStatus contains the status value of all modules.
var moduleStatus ModuleStatus

// RegisterModule adds a module to the global registry.
// Typically called from an init() function inside each module package.
func RegisterModule(module Module) {
	moduleRegistry = append(moduleRegistry, module)
	moduleStatus.ModuleCount = moduleStatus.ModuleCount + 1
}

// LoadRegisteredModules loads and starts all modules that registered themselves via init()
func loadRegisteredModules(mux *http.ServeMux, ignoreErrors bool) error {
	for _, module := range moduleRegistry {
		// attempt to load module config
		config, okConfig := GetModuleConfig(module.Name)
		if !okConfig {
			// Failed to load module config
			moduleStatus.FailedModules = append(moduleStatus.FailedModules, module.Name)

			fErr := fmt.Sprintf("[MODULE] %s has no config!", module.Name)
			log.Println(fErr)

			if !ignoreErrors {
				return errors.New(fErr)
			}

			continue
		}

		// get module config enabled
		enable, okEnable := cfg.GetNestedValue[bool](config, "enable")
		if !okEnable {
			// The config has no enable value.
			moduleStatus.FailedModules = append(moduleStatus.FailedModules, module.Name)

			fErr := fmt.Sprintf("[MODULE] %s has no enable value in config!", module.Name)
			log.Println(fErr)

			if !ignoreErrors {
				return errors.New(fErr)
			}

			continue
		}

		// Skip further processing if module is disabled
		if !enable {
			// Module is disabled
			moduleStatus.DisabledModules = append(moduleStatus.DisabledModules, module.Name)
			log.Printf("[MODULE] %s is disabled.", module.Name)
			continue
		}

		// Check if module has Routes() set
		if module.Routes == nil {
			moduleStatus.FailedModules = append(moduleStatus.FailedModules, module.Name)
			fErr := fmt.Sprintf("[MODULE] %s has no routes!", module.Name)
			log.Println(fErr)

			if !ignoreErrors {
				return errors.New(fErr)
			}

			continue
		}

		// Try to start module
		startErr := safeStartModule(module)
		if startErr != nil {
			// Error or panic occured while trying to start module
			moduleStatus.FailedModules = append(moduleStatus.FailedModules, module.Name)
			log.Println(startErr)

			if !ignoreErrors {
				return startErr
			}

			continue
		}

		// Module successfully started, registering now.
		// Mount each module under its prefix.
		prefix, handler := module.Routes()
		mux.Handle(prefix, handler)
		moduleStatus.RunningModules = append(moduleStatus.RunningModules, module.Name)
	}

	return nil
}

// ShutdownRegisteredModules shuts down all modules that are running.
func shutdownRegisteredModules() {
	for _, module := range moduleRegistry {
		// Only shutdown module if it was started
		if slices.Contains(moduleStatus.RunningModules, module.Name) {
			err := safeShutdownModule(module)
			if err != nil {
				// Error or panic occured while trying to shutdown module
				log.Println(err)
				continue
			}
		}
	}
}

// safeShutdownModule is a function that attempts to execute the shutdown function of a module.
// It returns nil if the shutdown is successfull or the module does not provide a shutdown function.
// It makes sure that a panic in the shutdown function does not crash the server.
func safeShutdownModule(module Module) (err error) {
	// return immediately if shutdown is not needed.
	if module.Shutdown == nil {
		return nil
	}

	// recover panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[MODULE] %s panicked during shutdown: %v", module.Name, r)
		}
	}()

	if startupErr := module.Shutdown(); startupErr != nil {
		return fmt.Errorf("[MODULE] %s failed shutdown: %w", module.Name, startupErr)
	}

	return nil
}

// safeStartModule is a function that attempts to execute the startup function of a module.
// It returns nil if the startup is successfull or the module does not provide a startup function.
// It makes sure that a panic in the startup function does not crash the server.
func safeStartModule(module Module) (err error) {
	// return immediately if startup is not needed.
	if module.Startup == nil {
		return nil
	}

	// recover panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[MODULE] %s panicked during startup: %v", module.Name, r)
		}
	}()

	if startupErr := module.Startup(); startupErr != nil {
		return fmt.Errorf("[MODULE] %s failed startup: %w", module.Name, startupErr)
	}

	return nil
}

// GetModuleStatus queries the moduleRegistry and returns a ModuleStatus of all registered modules.
// This can be used by external monitoring tools or by modules (like builtins/health)
func GetModuleStatus() ModuleStatus {
	return moduleStatus
}
