package health

import (
	"encoding/json"
	"net/http"

	"github.com/karotte128/karotteapi/v2/api"
)

func health(w http.ResponseWriter, r *http.Request) {

	type response struct {
		ApiStatus         string `json:"apiStatus"`
		TotalModules      int    `json:"totalModules"`
		RegisteredModules int    `json:"registeredModules"`
		RunningModules    int    `json:"runningModules"`
		DisabledModules   int    `json:"disabledModules"`
		FailedModules     int    `json:"failedModules"`
	}

	status := api.GetModuleStatus()

	var apiStatus string

	if status.ModuleCount == len(status.RunningModules)+len(status.DisabledModules) {
		apiStatus = "ok"
	} else {
		apiStatus = "degraded"
	}

	req_response := response{
		ApiStatus:         apiStatus,
		TotalModules:      status.ModuleCount,
		RegisteredModules: len(status.RegisteredModules),
		RunningModules:    len(status.RunningModules),
		DisabledModules:   len(status.DisabledModules),
		FailedModules:     len(status.FailedModules),
	}

	json.NewEncoder(w).Encode(req_response)
}
