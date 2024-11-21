package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mwantia/nautilus/internal/agent"
	"github.com/mwantia/nautilus/plugins/debug"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nautilus",
		Short: "Nautilus document processing system",
	}

	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Run the Nautilus agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := agent.NewAgent()
			return a.Serve()
		},
	}

	var pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			network, _ := cmd.Flags().GetString("network")
			address, _ := cmd.Flags().GetString("address")

			if strings.TrimSpace(network) != "" && strings.TrimSpace(address) != "" {
				switch name {
				case "debug":
					debug.ServeOverNetwork(network, address)
				}
			} else {
				switch name {
				case "debug":
					return debug.Serve()
				}
			}

			return fmt.Errorf("unknown plugin: %s", name)
		},
	}

	pluginCmd.Flags().String("name", "", "Plugin name")
	pluginCmd.Flags().String("network", "tcp", "Plugin network")
	pluginCmd.Flags().String("address", "", "Plugin address")

	rootCmd.AddCommand(agentCmd, pluginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
