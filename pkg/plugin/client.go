package plugin

import (
	"net/rpc"
)

type RpcClient struct {
	Client *rpc.Client
}

func (c *RpcClient) Name() (string, error) {
	var resp string
	if err := c.Client.Call("Plugin.Name", struct{}{}, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

func (c *RpcClient) GetCapabilities() (PipelineProcessorCapability, error) {
	var resp PipelineProcessorCapability
	if err := c.Client.Call("Plugin.GetCapabilities", struct{}{}, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *RpcClient) Configure(cfg map[string]interface{}) error {
	var resp error
	err := c.Client.Call("Plugin.Configure", cfg, &resp)
	return err
}

func (c *RpcClient) Process(data *PipelineContextData) (*PipelineContextData, error) {
	var resp PipelineContextData
	err := c.Client.Call("Plugin.Process", data, &resp)
	return &resp, err
}

func (c *RpcClient) Health() error {
	var resp error
	err := c.Client.Call("Plugin.Health", struct{}{}, &resp)
	return err
}
