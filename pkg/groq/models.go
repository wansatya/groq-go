package groq

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Model represents information about a language model
type Model struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    OwnedBy string `json:"owned_by"`
}

// ModelList represents a list of models
type ModelList struct {
    Object string  `json:"object"`
    Data   []Model `json:"data"`
}

type modelCache struct {
    models   map[string]bool
    lastFetch time.Time
    mutex    sync.RWMutex
}

var (
    cache modelCache
    cacheDuration = 1 * time.Hour
)

// ListModels retrieves a list of all available models
func (c *Client) ListModels(ctx context.Context) (*ModelList, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/models", nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+c.APIKey)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
    }

    var modelList ModelList
    if err := json.NewDecoder(resp.Body).Decode(&modelList); err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    // Update the cache
    cache.mutex.Lock()
    cache.models = make(map[string]bool)
    for _, model := range modelList.Data {
        cache.models[model.ID] = true
    }
    cache.lastFetch = time.Now()
    cache.mutex.Unlock()

    return &modelList, nil
}

// GetModel retrieves information about a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/models/"+modelID, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+c.APIKey)

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
    }

    var model Model
    if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    return &model, nil
}

// IsValidModel checks if the given model name is valid
func (c *Client) IsValidModel(ctx context.Context, modelID string) (bool, error) {
    cache.mutex.RLock()
    if time.Since(cache.lastFetch) < cacheDuration {
        isValid := cache.models[modelID]
        cache.mutex.RUnlock()
        return isValid, nil
    }
    cache.mutex.RUnlock()

    // If cache is expired or empty, refresh it
    _, err := c.ListModels(ctx)
    if err != nil {
        return false, fmt.Errorf("error refreshing model list: %w", err)
    }

    cache.mutex.RLock()
    isValid := cache.models[modelID]
    cache.mutex.RUnlock()

    return isValid, nil
}