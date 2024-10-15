package groq

import (
    "bufio"
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

// ChatCompletionRequest represents a request to the chat completions endpoint
type ChatCompletionRequest struct {
    Model          string          `json:"model"`
    Messages       []Message       `json:"messages"`
    MaxTokens      int             `json:"max_tokens,omitempty"`
    Temperature    float32         `json:"temperature,omitempty"`
    TopP           float32         `json:"top_p,omitempty"`
    N              int             `json:"n,omitempty"`
    Stream         bool            `json:"stream,omitempty"`
    Stop           []string        `json:"stop,omitempty"`
    PresencePenalty float32        `json:"presence_penalty,omitempty"`
    FrequencyPenalty float32       `json:"frequency_penalty,omitempty"`
    LogitBias      map[string]int  `json:"logit_bias,omitempty"`
    User           string          `json:"user,omitempty"`
    ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// Message represents a single message in a chat completion request or response
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ResponseFormat represents the format of the response
type ResponseFormat struct {
    Type string `json:"type,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completions endpoint
type ChatCompletionResponse struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`
    Choices []struct {
        Index        int     `json:"index"`
        Message      Message `json:"message"`
        FinishReason string  `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

// ChatCompletionChunk represents a chunk of a streaming response
type ChatCompletionChunk struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`
    Choices []struct {
        Index        int     `json:"index"`
        Delta        Message `json:"delta"`
        FinishReason string  `json:"finish_reason"`
    } `json:"choices"`
}

// CreateChatCompletion sends a request to the chat completions endpoint
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
    if req.Stream {
        return nil, fmt.Errorf("streaming is not supported in this method, use CreateChatCompletionStream instead")
    }

    isValid, err := c.IsValidModel(ctx, req.Model)
    if err != nil {
        return nil, fmt.Errorf("error validating model: %w", err)
    }
    if !isValid {
        return nil, fmt.Errorf("invalid model: %s", req.Model)
    }

    // Include system prompts if set
    if len(c.SystemPrompts) > 0 {
        req.Messages = append(c.SystemPrompts, req.Messages...)
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

    body, err := io.ReadAll(resp.Body)
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

    // Handle potential JSON response format
    if req.ResponseFormat != nil && req.ResponseFormat.Type == "json_object" {
        // Attempt to parse the content as JSON
        var jsonContent interface{}
        err = json.Unmarshal([]byte(result.Choices[0].Message.Content), &jsonContent)
        if err != nil {
            return nil, fmt.Errorf("error parsing JSON content: %w", err)
        }
        // Re-encode the parsed JSON to ensure it's properly formatted
        formattedJSON, err := json.MarshalIndent(jsonContent, "", "  ")
        if err != nil {
            return nil, fmt.Errorf("error formatting JSON content: %w", err)
        }
        result.Choices[0].Message.Content = string(formattedJSON)
    }

    return &result, nil
}

// CreateChatCompletionStream sends a streaming request to the chat completions endpoint
func (c *Client) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan ChatCompletionChunk, <-chan error) {
    chunkChan := make(chan ChatCompletionChunk)
    errChan := make(chan error, 1)

    go func() {
        defer close(chunkChan)
        defer close(errChan)

        isValid, err := c.IsValidModel(ctx, req.Model)
        if err != nil {
            errChan <- fmt.Errorf("error validating model: %w", err)
            return
        }
        if !isValid {
            errChan <- fmt.Errorf("invalid model: %s", req.Model)
            return
        }

        // Include system prompts if set
        if len(c.SystemPrompts) > 0 {
            req.Messages = append(c.SystemPrompts, req.Messages...)
        }

        // Ensure stream is true for streaming requests
        req.Stream = true

        jsonData, err := json.Marshal(req)
        if err != nil {
            errChan <- fmt.Errorf("error marshaling request: %w", err)
            return
        }

        httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
        if err != nil {
            errChan <- fmt.Errorf("error creating request: %w", err)
            return
        }

        httpReq.Header.Set("Content-Type", "application/json")
        httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

        resp, err := c.HTTPClient.Do(httpReq)
        if err != nil {
            errChan <- fmt.Errorf("error sending request: %w", err)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            errChan <- fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
            return
        }

        reader := bufio.NewReader(resp.Body)
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err == io.EOF {
                    return
                }
                errChan <- fmt.Errorf("error reading stream: %w", err)
                return
            }

            line = strings.TrimSpace(line)
            if line == "" {
                continue
            }
            if !strings.HasPrefix(line, "data: ") {
                continue
            }
            line = strings.TrimPrefix(line, "data: ")
            if line == "[DONE]" {
                return
            }

            var chunk ChatCompletionChunk
            err = json.Unmarshal([]byte(line), &chunk)
            if err != nil {
                errChan <- fmt.Errorf("error unmarshaling chunk: %w", err)
                return
            }

            chunkChan <- chunk
        }
    }()

    return chunkChan, errChan
}