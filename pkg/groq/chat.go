package groq

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// ResponseFormat represents the format of the response
type ResponseFormat struct {
    Type string `json:"type,omitempty"`
}

// ChatCompletionRequest represents a request to the chat completions endpoint
type ChatCompletionRequest struct {
    Model          string          `json:"model"`
    Messages       []Message       `json:"messages"`
    MaxTokens      int             `json:"max_tokens,omitempty"`
    Temperature    float32         `json:"temperature,omitempty"`
    ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// Message represents a single message in a chat completion request or response
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ChatCompletionResponse represents a response from the chat completions endpoint
type ChatCompletionResponse struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`
    Choices []struct {
        Index   int     `json:"index"`
        Message Message `json:"message"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

// CreateChatCompletion sends a request to the chat completions endpoint
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
    if !IsValidModel(req.Model) {
        return nil, fmt.Errorf("invalid model: %s", req.Model)
    }

    // Include system prompts if set
    if len(c.SystemPrompts) > 0 {
        req.Messages = append(c.SystemPrompts, req.Messages...)
    }

    // If ResponseFormat is not set, it will default to text
    if req.ResponseFormat != nil && req.ResponseFormat.Type != "json_object" {
        return nil, fmt.Errorf("invalid response format: %s", req.ResponseFormat.Type)
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("error marshaling request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

    resp, err := c.HTTPClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
    }

    var result ChatCompletionResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("error unmarshaling response: %w", err)
    }

    return &result, nil
}