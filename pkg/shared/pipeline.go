package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type PipelinePlugin struct {
	Impl PipelineProcessor
}

func (p *PipelinePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: p.Impl,
	}, nil
}

func (p *PipelinePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
	}, nil
}
