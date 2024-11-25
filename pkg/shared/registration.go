package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hashicorp/go-plugin"
)

type ServeRegistrationServer struct {
	Listener          net.Listener
	Server            *rpc.Server
	WaitGroup         sync.WaitGroup
	ActiveConnections sync.Map
	Cancel            context.CancelFunc
}

type ServeRegistrationErrorResponse struct {
	Error string `json:"error"`
}

func NewRegistrationServer(listener net.Listener) *ServeRegistrationServer {
	return &ServeRegistrationServer{
		Listener: listener,
		Server:   rpc.NewServer(),
	}
}

func (s *ServeRegistrationServer) Start(impl PipelineProcessor) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	s.Cancel = cancel
	if err := s.Server.RegisterName("Plugin", &RpcServer{Impl: impl}); err != nil {
		return err
	}

	s.WaitGroup.Add(1)
	go func() {
		defer s.WaitGroup.Done()
		s.AcceptConnection(ctx)
	}()

	return nil
}

func (s *ServeRegistrationServer) AcceptConnection(ctx context.Context) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return

			default:
				if errors.Is(err, net.ErrClosed) {
					return // Listener was closed, expected error
				}

				if ctx.Err() != nil {
					return // Context cancelled, expected error
				}

				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}

		s.WaitGroup.Add(1)
		go func(conn net.Conn) {
			defer s.WaitGroup.Done()
			defer conn.Close()

			s.ActiveConnections.Store(conn, struct{}{})
			defer s.ActiveConnections.Delete(conn)

			s.Server.ServeConn(conn)
		}(conn)
	}
}

func (s *ServeRegistrationServer) WaitForShutdown() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	if err := s.Listener.Close(); err != nil {
		return err
	}
	s.Cancel()

	s.ActiveConnections.Range(func(key, value interface{}) bool {
		if conn, ok := key.(net.Conn); ok {
			conn.Close()
		}
		return true
	})

	s.WaitGroup.Wait()
	return nil
}

func ServePipelineProcessor(impl PipelineProcessor) error {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"pipeline": &PipelinePlugin{
				Impl: impl,
			},
		},
	})

	return nil
}

func RegisterPipelineProcessor(impl PipelineProcessor, address string) error {
	name, err := impl.Name()
	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{
		"name":    name,
		"type":    "tcp",
		"address": "127.0.0.1:12345",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		return fmt.Errorf("failed to listen to address '%s': %v", "127.0.0.1:12345", err)
	}
	server := NewRegistrationServer(listener)
	if err := server.Start(impl); err != nil {
		return fmt.Errorf("failed to start registration-server: %v", err)
	}

	buf := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/plugin/register", address)

	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		server.Cancel() // Try to close any connections

		return fmt.Errorf("failed to post service registration to '%s': %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		server.Cancel() // Try to close any connections

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		var resp ServeRegistrationErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %v", err)
		}

		return fmt.Errorf(resp.Error)
	}

	return server.WaitForShutdown()
}
