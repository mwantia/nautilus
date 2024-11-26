package agent

import (
	"net/http"

	"github.com/mwantia/nautilus/internal/handler"
)

func SetupServer(a *NautilusAgent) (*http.Server, error) {
	mux := http.NewServeMux()
	if err := AddRoutes(a, mux); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    a.Config.Agent.Address,
		Handler: mux,
	}

	return server, nil
}

func AddRoutes(a *NautilusAgent, mux *http.ServeMux) error {
	mux.HandleFunc("/health", handler.HandleHealth(a.Registry))

	if err := AddApiRoutes(a, mux); err != nil {
		return err
	}

	mux.Handle("/", http.NotFoundHandler())
	return nil
}

func AddApiRoutes(a *NautilusAgent, mux *http.ServeMux) error {
	mux.HandleFunc("/v1/plugin/list", handler.HandlListPlugins(a.Registry))
	mux.HandleFunc("/v1/plugin/register", handler.HandleRegisterPlugin(a.Registry, a.Config))

	return nil
}
