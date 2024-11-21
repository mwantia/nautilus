package core

import (
	"encoding/json"
	"net/http"

	"github.com/mwantia/nautilus/pkg/shared"
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

func Health(plugins map[string]shared.NautilusPipelineProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := HealthResult{
			Status: "OK",
		}
		healthy := true

		for name, plugin := range plugins {
			sr := HealthPluginResult{
				Name:   name,
				Status: "OK",
			}

			if err := plugin.Health(); err != nil {
				sr.Status = "ERROR"
				sr.Error = err.Error()
				healthy = false
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
