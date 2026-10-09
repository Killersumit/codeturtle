package provider

import (
	"fmt"
	"os"

	"github.com/killersumit/codeturtle/internal/config"
)

// New creates and configures a Provider based on the given configuration.
// It retrieves the API key from the environment variable specified in cfg.Provider.APIKeyEnv.
func New(cfg *config.Config) (Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	keyEnv := cfg.Provider.APIKeyEnv
	if keyEnv == "" {
		return nil, fmt.Errorf("provider api_key_env is not specified in configuration")
	}

	apiKey := os.Getenv(keyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("API key environment variable %q is not set or empty. Please set it (e.g. export %s=your_api_key)", keyEnv, keyEnv)
	}

	return NewOpenAIProvider(
		cfg.Provider.BaseURL,
		apiKey,
		cfg.Provider.Model,
		cfg.Provider.MaxTokens,
		cfg.Provider.Temperature,
	), nil
}
