package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/killersumit/codeturtle/internal/config"
)

func TestOpenAIProvider_Chat_Success(t *testing.T) {
	expectedReply := "Hello, I am Qwen."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", auth)
		}

		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: expectedReply,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenAIProvider(server.URL, "test-key", "test-model", 1024, 0.5)
	if p.Name() != "test-model" {
		t.Errorf("expected Name() %q, got %q", "test-model", p.Name())
	}

	reply, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "Hi"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != expectedReply {
		t.Errorf("expected reply %q, got %q", expectedReply, reply)
	}
}

func TestOpenAIProvider_Chat_Non200Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": {"message": "Invalid API key"}}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider(server.URL, "invalid-key", "test-model", 1024, 0.2)
	_, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "401") {
		t.Errorf("expected error to mention 401, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "Invalid API key") {
		t.Errorf("expected error to contain body snippet, got: %s", errMsg)
	}
}

func TestOpenAIProvider_Chat_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices": []}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider(server.URL, "key", "test-model", 1024, 0.2)
	_, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
	if !strings.Contains(err.Error(), "empty choices") {
		t.Errorf("expected error to mention empty choices, got: %v", err)
	}
}

func TestOpenAIProvider_Chat_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json at all`))
	}))
	defer server.Close()

	p := NewOpenAIProvider(server.URL, "key", "test-model", 1024, 0.2)
	_, err := p.Chat(context.Background(), []Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parsing response JSON") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestOpenAIProvider_Chat_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := NewOpenAIProvider(server.URL, "key", "test-model", 1024, 0.2)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := p.Chat(ctx, []Message{{Role: "user", Content: "Hi"}})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestFactory_New_MissingEnv(t *testing.T) {
	_ = os.Unsetenv("TEST_CODETURTLE_KEY_MISSING")

	cfg := &config.Config{
		Provider: config.ProviderConfig{
			BaseURL:   "https://openrouter.ai/api/v1",
			Model:     "qwen/qwen3.8-27b",
			APIKeyEnv: "TEST_CODETURTLE_KEY_MISSING",
		},
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error when API key env var is missing, got nil")
	}
	if !strings.Contains(err.Error(), "TEST_CODETURTLE_KEY_MISSING") {
		t.Errorf("expected error to mention env var name, got: %v", err)
	}
}

func TestFactory_New_Success(t *testing.T) {
	t.Setenv("TEST_CODETURTLE_KEY_PRESENT", "my-secret-key")

	cfg := &config.Config{
		Provider: config.ProviderConfig{
			BaseURL:   "https://openrouter.ai/api/v1",
			Model:     "qwen/qwen3.8-27b",
			APIKeyEnv: "TEST_CODETURTLE_KEY_PRESENT",
		},
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Provider")
	}
}
