package exporter

import (
	"time"
)

// Config holds the configuration for the BCOS metrics exporter tool.
type Config struct {
	// Execution is the BCOS node to use.
	Execution ExecutionNode `yaml:"execution"`
	// DiskUsage determines if the disk usage metrics should be exported.
	DiskUsage DiskUsage `yaml:"diskUsage"`
}

// ExecutionNode represents a single FISCO-BCOS node.
type ExecutionNode struct {
	Enabled bool     `yaml:"enabled"`
	Name    string   `yaml:"name"`
	URL     string   `yaml:"url"`
	Modules []string `yaml:"modules"`
}

// DiskUsage configures the exporter to expose disk usage stats for these directories.
type DiskUsage struct {
	Enabled     bool          `yaml:"enabled"`
	Directories []string      `yaml:"directories"`
	Interval    time.Duration `yaml:"interval"`
}

// DefaultConfig represents a sane-default configuration.
func DefaultConfig() *Config {
	return &Config{
		Execution: ExecutionNode{
			Enabled: true,
			Name:    "bcos-node",
			URL:     "http://localhost:8545",
			Modules: []string{"bcos", "net", "web3"},
		},
		DiskUsage: DiskUsage{
			Enabled:     false,
			Directories: []string{},
			Interval:    60 * time.Minute,
		},
	}
}
