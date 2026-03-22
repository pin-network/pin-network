// Package config handles loading and validating meshd configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the top-level meshd configuration structure.
type Config struct {
	Node     NodeConfig     `yaml:"node"`
	Schedule ScheduleConfig `yaml:"schedule"`
	Network  NetworkConfig  `yaml:"network"`
	Limits   LimitsConfig   `yaml:"limits"`
	DataDir  string         `yaml:"-"`
	Dev      bool           `yaml:"-"`
}

// NodeConfig defines what this node contributes to the network.
type NodeConfig struct {
	Tier           int    `yaml:"tier"`
	StoragePath    string `yaml:"storage_path"`
	StorageLimitGB int    `yaml:"storage_limit_gb"`
	BandwidthMbps  int    `yaml:"bandwidth_limit_mbps"`
}

// ScheduleConfig defines when this node is active.
type ScheduleConfig struct {
	AlwaysOn         bool          `yaml:"always_on"`
	ActiveHours      []ActiveHours `yaml:"active_hours"`
	HeavyTasksDuring string        `yaml:"heavy_tasks_only_during"`
}

// ActiveHours defines a time window for hosting.
type ActiveHours struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

// NetworkConfig defines network behaviour.
type NetworkConfig struct {
	ListenPort     int      `yaml:"listen_port"`
	APIPort        int      `yaml:"api_port"`
	EnableUPnP     bool     `yaml:"enable_upnp"`
	EnableRelay    bool     `yaml:"enable_relay"`
	BootstrapNodes []string `yaml:"bootstrap_nodes"`
	PeerAPIs       []string `yaml:"peer_apis"`
}

// LimitsConfig defines resource limits for this node.
type LimitsConfig struct {
	CPUPercent    int  `yaml:"cpu_percent"`
	RAMMB         int  `yaml:"ram_mb"`
	BandwidthMbps int  `yaml:"bandwidth_mbps"`
	BatteryMinPct int  `yaml:"battery_min_percent"`
	WiFiOnly      bool `yaml:"wifi_only"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %w", err)
		}
		path = filepath.Join(home, ".pin", "config.yaml")
	}

	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *Config) LedgerPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pin", "ledger.db")
}

func (c *Config) StorePath() string {
	if c.Node.StoragePath != "" {
		return c.Node.StoragePath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pin", "store")
}

func defaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Node: NodeConfig{
			Tier:           1,
			StoragePath:    filepath.Join(home, ".pin", "store"),
			StorageLimitGB: 10,
			BandwidthMbps:  5,
		},
		Schedule: ScheduleConfig{
			AlwaysOn: true,
		},
		Network: NetworkConfig{
			ListenPort:  4001,
			APIPort:     4002,
			EnableUPnP:  true,
			EnableRelay: true,
			BootstrapNodes: []string{
				"/dns4/bootstrap.pin.network/tcp/4001/p2p/QmPlaceholderBootstrapID",
			},
		},
		Limits: LimitsConfig{
			CPUPercent:    25,
			RAMMB:         256,
			BandwidthMbps: 5,
			BatteryMinPct: 30,
			WiFiOnly:      false,
		},
	}
}

func (c *Config) validate() error {
	if c.Node.Tier < 1 || c.Node.Tier > 3 {
		return fmt.Errorf("node.tier must be 1, 2, or 3 (got %d)", c.Node.Tier)
	}
	if c.Network.ListenPort < 1024 || c.Network.ListenPort > 65535 {
		return fmt.Errorf("network.listen_port must be between 1024 and 65535")
	}
	if c.Network.APIPort < 1024 || c.Network.APIPort > 65535 {
		return fmt.Errorf("network.api_port must be between 1024 and 65535")
	}
	if c.Limits.CPUPercent < 1 || c.Limits.CPUPercent > 100 {
		return fmt.Errorf("limits.cpu_percent must be between 1 and 100")
	}
	return nil
}
