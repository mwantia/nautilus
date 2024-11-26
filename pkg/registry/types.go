package registry

import (
	"sync"
	"time"

	"github.com/mwantia/nautilus/pkg/plugin"
)

type PluginCleanup func() error

type PluginInfo struct {
	Name           string                             `json:"name"`
	Address        string                             `json:"address"`
	IsNetwork      bool                               `json:"is_network"`
	LastSeen       time.Time                          `json:"last_seen"`
	LastKnownError error                              `json:"-"`
	IsHealthy      bool                               `json:"is_healthy"`
	Capabilities   plugin.PipelineProcessorCapability `json:"capabilities"`
	Processor      plugin.PipelineProcessor           `json:"-"`
	Cleanup        PluginCleanup                      `json:"-"`
}

type PluginRegistry struct {
	Mutex   sync.RWMutex
	Plugins map[string]*PluginInfo
}
