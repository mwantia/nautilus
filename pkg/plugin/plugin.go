package plugin

import (
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "NAUTLIUS",
	MagicCookieValue: "pipeline",
}

type PipelineContextData struct {
	ID          string
	ContentType string
	Data        []byte
	Metadata    map[string]interface{}
}

type PipelineProcessor interface {
	Name() (string, error)

	GetCapabilities() (PipelineProcessorCapability, error)

	Configure(cfg map[string]interface{}) error

	Process(data *PipelineContextData) (*PipelineContextData, error)

	Health() error
}

var Plugins = map[string]plugin.Plugin{
	"pipeline": &PipelinePlugin{},
}
