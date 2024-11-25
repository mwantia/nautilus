package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mwantia/nautilus/internal/agent"
	"github.com/mwantia/nautilus/internal/config"
	"github.com/mwantia/nautilus/pkg/shared"
	"github.com/mwantia/nautilus/plugins/debug"
	"github.com/spf13/cobra"
)

var (
	Config string
	Cfg    *config.NautilusConfig

	Root = &cobra.Command{
		Use:               "nautilus",
		Short:             "Nautilus document processing system",
		PersistentPreRunE: SetupConfig,
	}
	Agent = &cobra.Command{
		Use:   "agent",
		Short: "Run Nautilus agent",
		RunE:  RunServeAgent,
	}
	Plugin = &cobra.Command{
		Use:   "plugin [name]",
		Short: "Run embedded Nautilus plugin",
		Args:  cobra.MaximumNArgs(1),
		RunE:  RunServePlugin,
	}
)

func SetupConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.ParseConfig(Config)
	if err != nil {
		return fmt.Errorf("unable to complete config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("unable to validate config: %v", err)
	}

	log.Printf("%v", cfg.Agent.Address)

	Cfg = cfg
	return nil
}

func RunServeAgent(cmd *cobra.Command, args []string) error {
	return agent.NewAgent(Cfg).Serve()
}

func RunServePlugin(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("plugin name is required")
	}

	address, _ := cmd.Flags().GetString("address")

	switch strings.TrimSpace(args[0]) {
	case "debug":
		impl := debug.NewImpl()
		if strings.TrimSpace(address) != "" {
			return shared.RegisterPipelineProcessor(impl, address)
		}

		return shared.ServePipelineProcessor(impl)
	}

	return fmt.Errorf("unknown plugin: %s", args[0])
}

func main() {
	Root.PersistentFlags().StringVar(&Config, "config", "", "Defines the configuration path used by this application")
	Plugin.Flags().String("address", "", "If defined, registers the plugin in network mode and tries to connect to the external agent via 'address'.")

	Root.AddCommand(Agent, Plugin)

	if err := Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
