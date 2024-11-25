package config

import "fmt"

func (c *NautilusConfig) Validate() error {
	if c.Agent == nil {
		return fmt.Errorf("agent configuration block is required")
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
