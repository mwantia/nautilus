package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/mwantia/nautilus/internal/config"
	"github.com/mwantia/nautilus/pkg/registry"
	"github.com/mwantia/nautilus/pkg/shared"
)

type RegisterPluginRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func HandlListPlugins(reg *registry.PluginRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		plugins := reg.ListPlugins()

		w.WriteHeader(http.StatusOK)
		encoder.Encode(plugins)
	}
}

func HandleRegisterPlugin(reg *registry.PluginRegistry, cfg *config.NautilusConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			encoder.Encode(map[string]string{
				"error": "Method 'POST' not allowed",
			})

			return
		}

		var request RegisterPluginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("Invalid request body: %v", err),
			})

			return
		}

		client, err := rpc.Dial(request.Type, request.Address)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("Unable to connect to external plugin '%s': %v", request.Name, err),
			})

			return
		}

		processor := &shared.RpcClient{
			Client: client,
		}
		name, err := processor.Name()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("Unable to validate processor for plugin '%s': %v", request.Name, err),
			})

			client.Close() // Try to close any connections
			return
		}

		if strings.TrimSpace(name) != strings.TrimSpace(request.Name) {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("The requested name '%s' doesn't match with the provided plugin '%s'", request.Name, name),
			})

			client.Close() // Try to close any connections
			return
		}

		info := &registry.PluginInfo{
			Name:      request.Name,
			IsNetwork: true,
			Address:   request.Address,
			Processor: processor,
			Cleanup:   client.Close,
		}
		if err := reg.Register(info); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("Unable to register plugin '%s': %v", request.Name, err),
			})

			client.Close() // Try to close any connections
			return
		}

		pluginCfg, err := cfg.GetPluginConfigMap(name)
		if err != nil {
			log.Printf("Error loading plugin config: %v", err)
		}

		if err := processor.Configure(pluginCfg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{
				"error": fmt.Sprintf("Unable to configure plugin '%s': %v", request.Name, err),
			})

			client.Close() // Try to close any connections
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder.Encode(map[string]string{
			"status": "OK",
		})
	}
}
