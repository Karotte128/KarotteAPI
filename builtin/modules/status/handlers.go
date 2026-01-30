package status

import (
	"encoding/json"
	"net/http"

	"github.com/karotte128/karotteapi/internal"
)

func status(w http.ResponseWriter, r *http.Request) {

	type response struct {
		ApiStatus         string `json:"apiStatus"`
		TotalModules      int    `json:"totalModules"`
		RegisteredModules int    `json:"registeredModules"`
		RunningModules    int    `json:"runningModules"`
		DisabledModules   int    `json:"disabledModules"`
		FailedModules     int    `json:"failedModules"`
	}

	status := internal.GetModuleStatus()

	var apiStatus string

	if status.TotalModules == status.RunningModules+status.DisabledModules {
		apiStatus = "ok"
	} else {
		apiStatus = "degraded"
	}

	req_response := response{
		ApiStatus:         apiStatus,
		TotalModules:      status.TotalModules,
		RegisteredModules: status.RegisteredModules,
		RunningModules:    status.RunningModules,
		DisabledModules:   status.DisabledModules,
		FailedModules:     status.FailedModules,
	}

	json.NewEncoder(w).Encode(req_response)
}
