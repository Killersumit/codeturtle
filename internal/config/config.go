package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// DefaultConfigPath is the default location for the CodeTurtle configuration file.
const DefaultConfigPath = ".codeturtle/config.toml"

// Default configuration values.
const (
	DefaultMaxTokens   = 4096
	DefaultTemperature = 0.2
)

// Config represents the top-level CodeTurtle configuration.
type Config struct {
	Provider ProviderConfig `toml:"provider"`
}

// ProviderConfig represents LLM provider configuration.
type ProviderConfig struct {
	BaseURL     string  `toml:"base_url"`
	Model       string  `toml:"model"`
	APIKeyEnv   string  `toml:"api_key_env"`
	MaxTokens   int     `toml:"max_tokens"`
	Temperature float64 `toml:"temperature"`
}

// Load reads and parses a TOML config file from the specified path.
// If path is empty, DefaultConfigPath (".codeturtle/config.toml") is used.
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := &Config{
		Provider: ProviderConfig{
			MaxTokens:   DefaultMaxTokens,
			Temperature: DefaultTemperature,
		},
	}

	meta, err := toml.Decode(string(data), cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// If max_tokens was not explicitly set in the config, use default.
	if !meta.IsDefined("provider", "max_tokens") {
		cfg.Provider.MaxTokens = DefaultMaxTokens
	}

	// If temperature was not explicitly set in the config, use default.
	if !meta.IsDefined("provider", "temperature") {
		cfg.Provider.Temperature = DefaultTemperature
	}

	return cfg, nil
}
