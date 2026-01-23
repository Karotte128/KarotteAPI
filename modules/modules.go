package modules

// This file acts as a “module loader” for the entire system.
//
// In Go, packages must be imported for their init() functions to run.
// Since all modules register themselves via init(), we only need to
// import each module here once. This means:
//
//   - main.go never needs to be modified when adding a new module
//   - adding a module = add a new `_ "module/path"` import here
//   - each module self-registers via core.RegisterModule(...)
//
//
// How to add a new module:
// ------------------------
// Simply create a new directory under /modules/<name>/
// and add a blank import here:
//
//   import _ "karotte128-api/modules/newmodule"
//
//
// The core system will pick it up automatically.

import (
	_ "github.com/karotte128/karotteapi/modules/status"
)
