package groq

import (
    "net/http"
    "time"
)

const (
    defaultBaseURL = "https://api.groq.com/openai/v1"
    defaultTimeout = 30 * time.Second
)

// Client is the main struct for interacting with the Groq API
type Client struct {
    BaseURL       string
    APIKey        string
    HTTPClient    *http.Client
    SystemPrompts []Message
}

// NewClient creates a new Groq API client
func NewClient(apiKey string) *Client {
    return &Client{
        BaseURL: defaultBaseURL,
        APIKey:  apiKey,
        HTTPClient: &http.Client{
            Timeout: defaultTimeout,
        },
    }
}

// SetBasePrompt adds a system prompt to be included in all requests
func (c *Client) SetBasePrompt(content string) {
    c.SystemPrompts = append(c.SystemPrompts, Message{
        Role:    "system",
        Content: content,
    })
}

// ClearBasePrompts removes all system prompts
func (c *Client) ClearBasePrompts() {
    c.SystemPrompts = nil
}