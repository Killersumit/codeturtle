package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIProvider communicates with OpenAI-compatible chat completion APIs.
type OpenAIProvider struct {
	baseURL     string
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	client      *http.Client
}

// NewOpenAIProvider creates a new OpenAI-compatible provider.
func NewOpenAIProvider(baseURL, apiKey, model string, maxTokens int, temperature float64) *OpenAIProvider {
	return &OpenAIProvider{
		baseURL:     baseURL,
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		client:      &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns the model name configured for this provider.
func (p *OpenAIProvider) Name() string {
	return p.model
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Chat sends a series of messages to {base_url}/chat/completions and returns the assistant's reply.
func (p *OpenAIProvider) Chat(ctx context.Context, messages []Message) (string, error) {
	reqBody := chatRequest{
		Model:       p.model,
		Messages:    messages,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling chat request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/chat/completions", strings.TrimRight(p.baseURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBytes))
	if err != nil {
		return "", fmt.Errorf("creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	// Handle non-200 status code
	if resp.StatusCode != http.StatusOK {
		snippet := strings.TrimSpace(string(bodyBytes))
		if len(snippet) > 500 {
			snippet = snippet[:500] + "..."
		}
		return "", fmt.Errorf("provider request failed with HTTP %d (%s): %s", resp.StatusCode, resp.Status, snippet)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		snippet := strings.TrimSpace(string(bodyBytes))
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return "", fmt.Errorf("parsing response JSON: %w (response snippet: %s)", err, snippet)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("provider returned empty choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
