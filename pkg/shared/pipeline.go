package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type NautilusPipelinePlugin struct {
	Impl NautilusPipelineProcessor
}

func (p *NautilusPipelinePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &NautilusRPCServer{Impl: p.Impl}, nil
}

// Client returns an RPC client for this plugin type
func (p *NautilusPipelinePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &NautilusRPCClient{Client: c}, nil
}
