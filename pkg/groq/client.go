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
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
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

// SetBaseURL allows changing the base URL for the API
func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

// SetTimeout allows changing the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}