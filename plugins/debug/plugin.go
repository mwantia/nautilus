package debug

import (
	"log"

	"github.com/mwantia/nautilus/pkg/shared"
)

type DebugProcessor struct {
}

func NewImpl() *DebugProcessor {
	return &DebugProcessor{}
}

func (p *DebugProcessor) Name() (string, error) {
	return "debug", nil
}

func (p *DebugProcessor) Process(ctx *shared.PipelineContextData) (*shared.PipelineContextData, error) {
	log.Println("Processing debug plugin...")
	return ctx, nil
}

func (p *DebugProcessor) Configure(cfg map[string]interface{}) error {
	for k, v := range cfg {
		log.Printf("Key: %s, Value: %v", k, v)
	}

	log.Println("Configuring debug plugin...")
	return nil
}

func (p *DebugProcessor) Health() error {
	log.Println("Checking health for debug plugin...")
	return nil
}
