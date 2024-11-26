package config

import (
	"fmt"
	"strings"
)

func (c *NautilusConfig) Validate() error {
	if c.Agent == nil {
		return fmt.Errorf("agent configuration block is required")
	}
	if strings.TrimSpace(c.LogLevel) == "" {
		c.LogLevel = "INFO"
	}

	if err := c.Agent.Validate(); err != nil {
		return fmt.Errorf("invalid agent configuration: %w", err)
	}

	return nil
}

func (c *AgentConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}
