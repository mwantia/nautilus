package agent

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/nautilus/pkg/registry"
	"github.com/mwantia/nautilus/pkg/shared"
)

func (agent *NautilusAgent) EmbbedPlugin(name string) error {
	path, err := os.Executable()
	if err != nil {
		return nil
	}

	if err = agent.LocalPlugin(path, "plugin", name); err != nil {
		return fmt.Errorf("unable to load embbeded plugin '%s': %v", name, err)
	}

	return nil
}

func (a *NautilusAgent) LocalPlugin(path string, arg ...string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.Plugins,
		Cmd:             exec.Command(path, arg...),
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

	processor := raw.(shared.PipelineProcessor)

	name, err := processor.Name()
	if err != nil {
		return err
	}

	log.Printf("Loaded local plugin named '%s'", name)

	info := &registry.PluginInfo{
		Name:      name,
		IsNetwork: false,
		Processor: processor,
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
		log.Printf("Error loading plugin config: %v", err)
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

	processor := &shared.RpcClient{
		Client: rpc,
	}
	name, err := processor.Name()
	if err != nil {
		return err
	}

	log.Printf("Loaded network plugin named '%s'", name)

	info := &registry.PluginInfo{
		Name:      name,
		IsNetwork: false,
		Processor: processor,
		Cleanup:   rpc.Close,
	}
	if err := a.Registry.Register(info); err != nil {
		return err
	}

	cfg, err := a.Config.GetPluginConfigMap(name)
	if err != nil {
		log.Printf("Error loading plugin config: %v", err)
	}

	if err := processor.Configure(cfg); err != nil {
		return err
	}

	return nil
}
