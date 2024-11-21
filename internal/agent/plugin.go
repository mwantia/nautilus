package agent

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/nautilus/pkg/shared"
)

func (agent *NautilusAgent) EmbbedPlugin(name string) error {
	path, err := os.Executable()
	if err != nil {
		return nil
	}

	if err = agent.LocalPlugin(path, "plugin", "--name", name); err != nil {
		return fmt.Errorf("unable to embbeded plugin '%s': %v", name, err)
	}

	return nil
}

func (agent *NautilusAgent) LocalPlugin(path string, arg ...string) error {
	agent.Mutex.Lock()
	defer agent.Mutex.Unlock()

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

	plugin := raw.(shared.NautilusPipelineProcessor)

	name, err := plugin.Name()
	if err != nil {
		return err
	}

	log.Printf("Loaded plugins named '%s'", name)

	if err := plugin.Configure(); err != nil {
		client.Kill()
		return err
	}

	agent.Clients[name] = client
	agent.Plugins[name] = plugin

	return nil
}

func (agent *NautilusAgent) NetworkPlugin(network, address string) error {
	agent.Mutex.Lock()
	defer agent.Mutex.Unlock()

	rpc, err := rpc.Dial(network, address)
	if err != nil {
		return err
	}

	plugin := &shared.NautilusRPCClient{
		Client: rpc,
	}
	name, err := plugin.Name()
	if err != nil {
		return err
	}

	if err := plugin.Configure(); err != nil {
		return err
	}

	agent.Plugins[name] = plugin

	return nil
}
