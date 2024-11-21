package shared

import "github.com/hashicorp/go-plugin"

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "NAUTLIUS",
	MagicCookieValue: "pipeline",
}

type NautilusPipelineContext struct {
	ID          string
	ContentType string
	Data        []byte
	Metadata    map[string]interface{}
}

type NautilusPipelineProcessor interface {
	Name() (string, error)

	Process(ctx *NautilusPipelineContext) (*NautilusPipelineContext, error)

	Configure() error

	Health() error
}

var Plugins = map[string]plugin.Plugin{
	"pipeline": &NautilusPipelinePlugin{},
}
