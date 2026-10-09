package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	tomlContent := `
[provider]
base_url = "https://openrouter.ai/api/v1"
model = "qwen/qwen3.8-27b"
api_key_env = "CODETURTLE_API_KEY"
`
	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Provider.BaseURL != "https://openrouter.ai/api/v1" {
		t.Errorf("expected BaseURL %q, got %q", "https://openrouter.ai/api/v1", cfg.Provider.BaseURL)
	}
	if cfg.Provider.Model != "qwen/qwen3.8-27b" {
		t.Errorf("expected Model %q, got %q", "qwen/qwen3.8-27b", cfg.Provider.Model)
	}
	if cfg.Provider.APIKeyEnv != "CODETURTLE_API_KEY" {
		t.Errorf("expected APIKeyEnv %q, got %q", "CODETURTLE_API_KEY", cfg.Provider.APIKeyEnv)
	}
	if cfg.Provider.MaxTokens != DefaultMaxTokens {
		t.Errorf("expected default MaxTokens %d, got %d", DefaultMaxTokens, cfg.Provider.MaxTokens)
	}
	if cfg.Provider.Temperature != DefaultTemperature {
		t.Errorf("expected default Temperature %f, got %f", DefaultTemperature, cfg.Provider.Temperature)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	tomlContent := `
[provider]
base_url = "https://api.openai.com/v1"
model = "gpt-4o"
api_key_env = "OPENAI_API_KEY"
max_tokens = 8192
temperature = 0.7
`
	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Provider.MaxTokens != 8192 {
		t.Errorf("expected MaxTokens 8192, got %d", cfg.Provider.MaxTokens)
	}
	if cfg.Provider.Temperature != 0.7 {
		t.Errorf("expected Temperature 0.7, got %f", cfg.Provider.Temperature)
	}
}

func TestLoad_ExplicitZeroTemperature(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")

	tomlContent := `
[provider]
base_url = "https://api.openai.com/v1"
model = "gpt-4o"
api_key_env = "OPENAI_API_KEY"
temperature = 0.0
`
	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Provider.Temperature != 0.0 {
		t.Errorf("expected Temperature 0.0, got %f", cfg.Provider.Temperature)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/path/that/does/not/exist.toml")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.toml")

	if err := os.WriteFile(configPath, []byte("invalid toml [[[" ), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid TOML, got nil")
	}
}
