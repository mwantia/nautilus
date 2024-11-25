package config

import "github.com/hashicorp/hcl/v2"

type NautilusConfig struct {
	Agent   *AgentConfig    `hcl:"agent,block"`
	Plugins []*PluginConfig `hcl:"plugin,block"`
}

type AgentConfig struct {
	Address   string `hcl:"address,optional"`
	PluginDir string `hcl:"plugin_dir,optional"`
}

type PluginConfig struct {
	Name    string           `hcl:"name,label"`
	Enabled bool             `hcl:"enabled,optional"`
	Config  PluginConfigBody `hcl:"config,block"`
}

type PluginConfigBody struct {
	Body hcl.Body `hcl:",remain"`
}

func NewDefault() *NautilusConfig {
	return &NautilusConfig{
		Agent: &AgentConfig{
			Address:   ":8080",
			PluginDir: "./plugins",
		},
	}
}
