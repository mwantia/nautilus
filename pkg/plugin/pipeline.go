package plugin

import (
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"
)

type PipelinePlugin struct {
	Impl PipelineProcessor
}

func (p *PipelinePlugin) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RpcServer{
		Impl: p.Impl,
	}, nil
}

func (p *PipelinePlugin) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{
		Client: c,
	}, nil
}
