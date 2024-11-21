package agent

import (
	"log"
	"net/http"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/nautilus/pkg/core"
	"github.com/mwantia/nautilus/pkg/shared"
)

type NautilusAgent struct {
	Mutex   sync.RWMutex
	Clients map[string]*plugin.Client
	Plugins map[string]shared.NautilusPipelineProcessor
	Config  NautilusConfig
}

func NewAgent() *NautilusAgent {
	return &NautilusAgent{
		Clients: make(map[string]*plugin.Client),
		Plugins: make(map[string]shared.NautilusPipelineProcessor),
		Config:  NautilusConfig{},
	}
}

func (agent *NautilusAgent) Serve() error {
	if err := agent.NetworkPlugin("tcp", "127.0.0.1:12345"); err != nil {
		log.Printf("Error loading plugin: %v", err)
	}

	http.HandleFunc("/health", core.Health(agent.Plugins))

	log.Println("Serving HTTP to :8080")
	return http.ListenAndServe(":8080", nil)
}

func (agent *NautilusAgent) Cleanup() {
	agent.Mutex.Lock()
	defer agent.Mutex.Unlock()

	for _, client := range agent.Clients {
		client.Kill()
	}
}
