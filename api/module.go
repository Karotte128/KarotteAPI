package api

import (
	"log"
	"net/http"

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
	// This can be used to cleanly disconnect from services connected during Startup().
	Shutdown func() error
}

// This is the module status type.
type status int

// status can have the following values:
const (
	statusRegistered status = iota
	statusRunning
	statusDisabled
	statusFailed
)

// registryModule is the internal data structure for the module registry.
// It is not public to the modules or the main package. It is only for use in core.
type registryModule struct {
	// module contains the data provided by the registering module
	module Module

	// status is the current status of the module.
	status status
}

// registry holds all globally registered modules.
var moduleRegistry []registryModule

// RegisterModule adds a module to the global registry.
// Typically called from an init() function inside each module package.
func RegisterModule(module Module) {
	// Structured data for the module registry
	var reg_mod = registryModule{
		module: module,
		status: statusRegistered,
	}

	// add module to registry
	moduleRegistry = append(moduleRegistry, reg_mod)
}

// LoadRegisteredModules loads and starts all modules that registered themselves via init()
func loadRegisteredModules(mux *http.ServeMux) {

	// Register and start all modules in the moduleRegistry
	for i, reg_mod := range moduleRegistry {
		var modStatus status
		var enabled bool = false

		// get enable value from config
		config, okConfig := GetModuleConfig(reg_mod.module.Name)
		if okConfig {
			enable_conf, okEnable := cfg.GetNestedValue[bool](config, "enable")
			if okEnable {
				enabled = enable_conf
			} else {
				// The config has no enable value.
				log.Printf("[MODULE] %s has no enable value in config!", reg_mod.module.Name)
			}
		} else {
			// The module has no config entry
			log.Printf("[MODULE] %s has no config!", reg_mod.module.Name)
		}

		// Test if module is enabled
		if enabled {
			// Module is running
			if reg_mod.module.Routes != nil {
				// Try to start module
				ok := safeStartModule(reg_mod.module)

				if ok {
					// Module successfully started, registering now.

					// Mount each module under its prefix.
					prefix, handler := reg_mod.module.Routes()
					mux.Handle(prefix, handler)

					// Set module status to running
					modStatus = statusRunning
				} else {
					// Module failed startup!
					log.Printf("[MODULE] %s was not registered!", reg_mod.module.Name)

					// Set module status to failed
					modStatus = statusFailed
				}
			} else {
				log.Printf("[MODULE] %s has no routes!", reg_mod.module.Name)

				// Set module status to failed
				modStatus = statusFailed
			}
		} else {
			// Module is disabled
			log.Printf("[MODULE] %s is disabled.", reg_mod.module.Name)

			// Set module status to disabled
			modStatus = statusDisabled
		}

		moduleRegistry[i].status = modStatus
	}
}

// ShutdownRegisteredModules shuts down all modules that are running.
func shutdownRegisteredModules() {
	for _, reg_mod := range moduleRegistry {
		if reg_mod.status == statusRunning {
			safeShutdownModule(reg_mod.module)
		}
	}
}

// safeShutdownModule is a function that attempts to execute the shutdown function of a module.
// It makes sure that a panic in the shutdown function does not crash the server.
func safeShutdownModule(module Module) {
	// only execute if the module implements a shutdown function
	if module.Shutdown != nil {
		// recover from panic
		defer func() {
			r := recover()
			if r != nil {
				log.Printf("[MODULE] %s panicked during shutdown: %v", module.Name, r)
			}
		}()

		// try to shutdown the module
		err := module.Shutdown()
		if err != nil {
			log.Printf("[MODULE] %s failed shutdown: %v", module.Name, err)
		}
	}
}

// safeStartModule is a function that attempts to execute the startup function of a module.
// It returns true if the startup is successfull or the module does not provide a startup function.
// It makes sure that a panic in the startup function does not crash the server.
func safeStartModule(module Module) bool {
	var ok bool = true

	// only execute if the module implements a startup function
	if module.Startup != nil {
		// recover from panic
		defer func() {
			r := recover()
			if r != nil {
				log.Printf("[MODULE] %s panicked during startup: %v", module.Name, r)
				ok = false
			}
		}()

		// try to start the module
		err := module.Startup()
		if err != nil {
			log.Printf("[MODULE] %s failed startup: %v", module.Name, err)
			ok = false
		}

		// returns true if startup was successfull
		return ok
	} else {
		// return true if startup is not needed
		return true
	}
}

// ModuleStatus holds status information about all registered modules.
// It contains lists of modules, sorted by status type.
type ModuleStatus struct {
	ModuleCount       int
	RegisteredModules []string
	RunningModules    []string
	DisabledModules   []string
	FailedModules     []string
}

// GetModuleStatus queries the moduleRegistry and returns a ModuleStatus of all registered modules.
// This can be used by external monitoring tools or by modules (like builtins/health)
func GetModuleStatus() ModuleStatus {
	var status ModuleStatus

	for _, module := range moduleRegistry {
		switch module.status {
		case statusRegistered:
			status.RegisteredModules = append(status.RegisteredModules, module.module.Name)

		case statusRunning:
			status.RunningModules = append(status.RunningModules, module.module.Name)

		case statusDisabled:
			status.DisabledModules = append(status.DisabledModules, module.module.Name)

		case statusFailed:
			status.FailedModules = append(status.FailedModules, module.module.Name)
		}
	}

	// ModuleCount is the total amount of registered modules.
	status.ModuleCount = len(moduleRegistry)

	return status
}
