package debug

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/mwantia/nautilus/pkg/shared"
)

type DebugProcessor struct {
}

func Serve() error {
	log.Println("Serving debug plugin...")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"pipeline": &shared.NautilusPipelinePlugin{
				Impl: &DebugProcessor{},
			},
		},
	})
	return nil
}

func ServeOverNetwork(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	defer listener.Close()

	fmt.Printf("Plugin server listening on %s://%s\n", network, address)

	server := rpc.NewServer()
	err = server.RegisterName("Plugin", &shared.NautilusRPCServer{
		Impl: &DebugProcessor{},
	})
	if err != nil {
		return fmt.Errorf("failed to register RPC server: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}
		go server.ServeConn(conn)
	}
}

func (p *DebugProcessor) Name() (string, error) {
	return "debug", nil
}

func (p *DebugProcessor) Process(ctx *shared.NautilusPipelineContext) (*shared.NautilusPipelineContext, error) {
	log.Println("Processing debug plugin...")
	return ctx, nil
}

func (p *DebugProcessor) Configure() error {
	log.Println("Configuring debug plugin...")
	return nil
}

func (p *DebugProcessor) Health() error {
	log.Println("Checking health for debug plugin...")
	return nil
}
