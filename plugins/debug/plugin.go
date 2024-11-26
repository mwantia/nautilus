package debug

import (
	"log"

	"github.com/mwantia/nautilus/pkg/plugin"
)

type DebugProcessor struct {
}

func NewImpl() *DebugProcessor {
	return &DebugProcessor{}
}

func (p *DebugProcessor) Name() (string, error) {
	return "debug", nil
}

func (p *DebugProcessor) GetCapabilities() (plugin.PipelineProcessorCapability, error) {
	return plugin.PipelineProcessorCapability{
		Types: []plugin.PipelineProcessorCapabilityType{
			plugin.None,
		},
	}, nil
}

func (p *DebugProcessor) Configure(cfg map[string]interface{}) error {
	for k, v := range cfg {
		log.Printf("Key: %s, Value: %v", k, v)
	}

	log.Println("Configuring debug plugin...")
	return nil
}

func (p *DebugProcessor) Process(data *plugin.PipelineContextData) (*plugin.PipelineContextData, error) {
	log.Println("Processing debug plugin...")
	return data, nil
}

func (p *DebugProcessor) Health() error {
	log.Println("Checking health for debug plugin...")
	return nil
}
