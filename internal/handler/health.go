package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mwantia/nautilus/pkg/registry"
)

type HealthResult struct {
	Status  string               `json:"status"`
	Plugins []HealthPluginResult `json:"plugins"`
}

type HealthPluginResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func HandleHealth(reg *registry.PluginRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := HealthResult{
			Status: "OK",
		}
		healthy := true

		for _, plugin := range reg.ListPlugins() {
			sr := HealthPluginResult{
				Name:   plugin.Name,
				Status: "OK",
			}

			if !plugin.IsHealthy {
				sr.Status = "ERROR"
				healthy = false

				if plugin.LastKnownError != nil {
					sr.Error = plugin.LastKnownError.Error()
				}
			}

			result.Plugins = append(result.Plugins, sr)
		}

		w.Header().Set("Content-Type", "application/json")

		if !healthy {
			result.Status = "ERROR"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		if err := encoder.Encode(result); err != nil {
			encoder.Encode(map[string]string{
				"error": err.Error(),
			})
		}
	}
}
