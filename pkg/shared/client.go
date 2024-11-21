package shared

import (
	"net/rpc"
)

type NautilusRPCClient struct {
	Client *rpc.Client
}

func (c *NautilusRPCClient) Name() (string, error) {
	var resp string
	if err := c.Client.Call("Plugin.Name", struct{}{}, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

func (c *NautilusRPCClient) Process(ctx *NautilusPipelineContext) (*NautilusPipelineContext, error) {
	var resp NautilusPipelineContext
	err := c.Client.Call("Plugin.Process", ctx, &resp)
	return &resp, err
}

func (c *NautilusRPCClient) Configure() error {
	var resp error
	err := c.Client.Call("Plugin.Configure", struct{}{}, &resp)
	return err
}

func (c *NautilusRPCClient) Health() error {
	var resp error
	err := c.Client.Call("Plugin.Health", struct{}{}, &resp)
	return err
}
