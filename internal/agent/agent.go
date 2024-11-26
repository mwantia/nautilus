package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mwantia/nautilus/internal/config"
	"github.com/mwantia/nautilus/pkg/log"
	"github.com/mwantia/nautilus/pkg/registry"
)

type NautilusAgent struct {
	Mutex    sync.RWMutex
	Logger   *log.Logger
	Registry *registry.PluginRegistry
	Config   *config.NautilusConfig
}

func NewAgent(cfg *config.NautilusConfig) *NautilusAgent {
	return &NautilusAgent{
		Registry: registry.NewRegistry(),
		Logger:   log.NewLogger("agent"),
		Config:   cfg,
	}
}

func (a *NautilusAgent) Serve(ctx context.Context) error {
	if err := a.ServeLocalPlugins(); err != nil {
		a.Logger.Warn("Unable to serve local plugin", "error", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv, err := SetupServer(a)
	if err != nil {
		return err
	}

	go a.Registry.Watch(ctx)

	go func() {
		a.Logger.Info("Serving HTTP server", "address", a.Config.Agent.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("Error serving http server", "address", a.Config.Agent.Address, "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		a.Logger.Info("Shutting down agent...")

		shutdown := context.Background()
		shutdown, cncl := context.WithTimeout(shutdown, 10*time.Second)
		defer cncl()

		if err := srv.Shutdown(shutdown); err != nil {
			a.Logger.Error("Error shutting down http server", "error", err)
		}
	}()

	wg.Wait()
	return a.Cleanup()
}

func (a *NautilusAgent) Cleanup() error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	var err error

	for _, plugin := range a.Registry.ListPlugins() {
		if cleanupErr := plugin.Cleanup(); cleanupErr != nil {
			err = errors.Join(err, cleanupErr)
		}
	}

	return err
}

func (a *NautilusAgent) ServeLocalPlugins() error {
	files, err := os.ReadDir(a.Config.Agent.PluginDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", a.Config.Agent.PluginDir, file.Name())
			if err := a.LocalPlugin(path); err != nil {
				a.Logger.Warn("Unable to load local plugin", "path", path, "error", err)
			}
		}
	}

	return nil
}
