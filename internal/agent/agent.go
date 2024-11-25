package agent

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/mwantia/nautilus/internal/config"
	"github.com/mwantia/nautilus/internal/handler"
	"github.com/mwantia/nautilus/pkg/registry"
)

type NautilusAgent struct {
	Mutex    sync.RWMutex
	Registry *registry.PluginRegistry
	Config   *config.NautilusConfig
}

func NewAgent(cfg *config.NautilusConfig) *NautilusAgent {
	return &NautilusAgent{
		Registry: registry.NewRegistry(),
		Config:   cfg,
	}
}

func (a *NautilusAgent) Serve() error {
	if err := a.ServeLocalPlugins(); err != nil {
		log.Printf("%v", err)
	}

	http.HandleFunc("/health", handler.HandleHealth(a.Registry))
	http.HandleFunc("/plugin/list", handler.HandlListPlugins(a.Registry))
	http.HandleFunc("/plugin/register", handler.HandleRegisterPlugin(a.Registry, a.Config))

	log.Printf("Serving HTTP to %s", a.Config.Agent.Address)
	return http.ListenAndServe(a.Config.Agent.Address, nil)
}

func (a *NautilusAgent) Cleanup() {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	for _, plugin := range a.Registry.ListPlugins() {
		if err := plugin.Cleanup(); err != nil {
			log.Printf("Plugin cleanup error: %v", err)
		}
	}
}

func (a *NautilusAgent) ServeLocalPlugins() error {
	files, err := os.ReadDir(a.Config.Agent.PluginDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", a.Config.Agent.PluginDir, file.Name())
			if err := a.LocalPlugin(path); err != nil {
				log.Printf("Unable to load local plugin !! '%s': %v", path, err)
			}
		}
	}

	return nil
}
