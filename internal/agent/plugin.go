package agent

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mwantia/nautilus/pkg/log"
	"github.com/mwantia/nautilus/pkg/plugin"
	"github.com/mwantia/nautilus/pkg/registry"
)

func (a *NautilusAgent) EmbbedPlugin(name string) error {
	path, err := os.Executable()
	if err != nil {
		return nil
	}

	if err = a.LocalPlugin(path, "plugin", name); err != nil {
		return fmt.Errorf("unable to load embbeded plugin '%s': %v", name, err)
	}

	return nil
}

func (a *NautilusAgent) LocalPlugin(path string, arg ...string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.Plugins,
		Cmd:             exec.Command(path, arg...),
		Logger:          log.Default,
	})

	rpc, err := client.Client()
	if err != nil {
		client.Kill()
		return err
	}

	raw, err := rpc.Dispense("pipeline")
	if err != nil {
		client.Kill()
		return err
	}

	processor := raw.(plugin.PipelineProcessor)

	name, err := processor.Name()
	if err != nil {
		return err
	}

	cap, err := processor.GetCapabilities()
	if err != nil {
		return err
	}

	a.Logger.Info("Loaded local plugin", "name", name)

	info := &registry.PluginInfo{
		Name:         name,
		IsNetwork:    false,
		Processor:    processor,
		Capabilities: cap,

		Cleanup: func() error {
			client.Kill()
			return nil
		},
	}
	if err := a.Registry.Register(info); err != nil {
		return err
	}

	cfg, err := a.Config.GetPluginConfigMap(name)
	if err != nil {
		a.Logger.Warn("Unable to load plugin config", "name", name, "error", err)
	}

	if err := processor.Configure(cfg); err != nil {
		client.Kill()
		return err
	}

	return nil
}

func (a *NautilusAgent) NetworkPlugin(network, address string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	rpc, err := rpc.Dial(network, address)
	if err != nil {
		return err
	}

	processor := &plugin.RpcClient{
		Client: rpc,
	}
	name, err := processor.Name()
	if err != nil {
		return err
	}

	cap, err := processor.GetCapabilities()
	if err != nil {
		return err
	}

	a.Logger.Info("Loaded network plugin", "name", name)

	info := &registry.PluginInfo{
		Name:         name,
		IsNetwork:    false,
		Processor:    processor,
		Capabilities: cap,
		Cleanup:      rpc.Close,
	}
	if err := a.Registry.Register(info); err != nil {
		return err
	}

	cfg, err := a.Config.GetPluginConfigMap(name)
	if err != nil {
		a.Logger.Warn("Unable to load plugin config", "name", name, "error", err)
	}

	if err := processor.Configure(cfg); err != nil {
		return err
	}

	return nil
}
