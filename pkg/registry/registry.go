package registry

import (
	"context"
	"fmt"
	"time"
)

func NewRegistry() *PluginRegistry {
	return &PluginRegistry{
		Plugins: make(map[string]*PluginInfo),
	}
}

func (reg *PluginRegistry) Watch(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		plugins := reg.ListPlugins()
		for _, plugin := range plugins {
			if err := plugin.Processor.Health(); err != nil {
				plugin.IsHealthy = false
				plugin.LastKnownError = err

				continue
			}

			plugin.IsHealthy = true
			plugin.LastKnownError = nil
			plugin.LastSeen = time.Now()
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (reg *PluginRegistry) Register(info *PluginInfo) error {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	if _, exists := reg.Plugins[info.Name]; exists {
		return fmt.Errorf("a plugin with the name %s has already been registered", info.Name)
	}

	info.LastSeen = time.Now()
	reg.Plugins[info.Name] = info

	return nil
}

func (reg *PluginRegistry) Deregister(name string) (*PluginInfo, error) {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("a plugin with the name %s does not exist", name)
	}

	return plugin, nil
}

func (reg *PluginRegistry) GetPlugin(name string) (*PluginInfo, bool) {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugin, exists := reg.Plugins[name]
	return plugin, exists
}

func (reg *PluginRegistry) ListPlugins() []*PluginInfo {
	reg.Mutex.Lock()
	defer reg.Mutex.Unlock()

	plugins := make([]*PluginInfo, 0, len(reg.Plugins))
	for _, plugin := range reg.Plugins {
		plugins = append(plugins, plugin)
	}

	return plugins
}
