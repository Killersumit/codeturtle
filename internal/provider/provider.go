package provider

import "context"

// Provider defines the interface for model-agnostic LLM providers.
type Provider interface {
	Chat(ctx context.Context, messages []Message) (string, error)
	Name() string
}

// Message represents a single chat message exchanged with a Provider.
type Message struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}
