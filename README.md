# Groq Go SDK

This is an unofficial Go SDK for the [Groq LPU™ AI inference](https://groq.com/inference). 

It provides a simple and efficient way to interact with Groq's AI models using Go.

## Project Structure

```
groq-go/
├── .github/
│   ├── ISSUE_TEMPLATE.md
│   └── PULL_REQUEST_TEMPLATE.md
├── cmd/
│   └── example/
│       └── main.go
├── docs/
│   └── API.md
├── internal/
│   └── util/
│       └── validation.go
├── pkg/
│   └── groq/
│       ├── client.go
│       ├── models.go
│       ├── chat.go
│       └── errors.go
├── test/
│   └── integration/
│       └── client_test.go
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```

## Installation

To install the Groq Go SDK, use `go get`:

```bash
go get github.com/wansatya/groq-go
```

## Usage

Here's a quick example of how to use the SDK:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/wansatya/groq-go/pkg/groq"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file")
    }

    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        log.Fatal("GROQ_API_KEY not found in environment variables")
    }

    modelID := os.Getenv("GROQ_MODEL")
    if modelID == "" {
        modelID = "mixtral-8x7b-32768" // Default model if not specified
        fmt.Printf("GROQ_MODEL not set, using default: %s\n", modelID)
    }

    client := groq.NewClient(apiKey)

    // Set base prompts
    client.SetBasePrompt("You are a helpful assistant. Always be polite and concise.")
    client.SetBasePrompt("Provide examples when explaining concepts.")

    ctx := context.Background()

    // List models
    listModels(ctx, client)

    // Get model details
    getModel(ctx, client, modelID)

    // Create chat completion
    chatCompletionCreate(ctx, client, modelID)

    // Create chat completion with JSON response
    chatCompletionCreateJSON(ctx, client, modelID)

    // Create streaming chat completion
    chatCompletionCreateStream(ctx, client, modelID)
}

func listModels(ctx context.Context, client *groq.Client) {
    modelList, err := client.ListModels(ctx)
    if err != nil {
        log.Fatalf("Error listing models: %v", err)
    }

    fmt.Println("Available Models:")
    for _, model := range modelList.Data {
        fmt.Printf("- %s\n", model.ID)
    }
}

func getModel(ctx context.Context, client *groq.Client, modelID string) {
    model, err := client.GetModel(ctx, modelID)
    if err != nil {
        log.Fatalf("Error fetching model %s: %v", modelID, err)
    }

    fmt.Printf("\nModel Details for %s:\n", modelID)
    fmt.Printf("ID: %s\n", model.ID)
    fmt.Printf("Object: %s\n", model.Object)
    fmt.Printf("Created: %d\n", model.Created)
    fmt.Printf("Owned By: %s\n", model.OwnedBy)

    isValid, err := client.IsValidModel(ctx, modelID)
    if err != nil {
        log.Fatalf("Error checking model validity: %v", err)
    }
    fmt.Printf("Is model %s valid: %v\n", modelID, isValid)
}

func chatCompletionCreate(ctx context.Context, client *groq.Client, modelID string) {
    req := groq.ChatCompletionRequest{
        Model: modelID,
        Messages: []groq.Message{
            {Role: "user", Content: "What is Golang?"},
        },
        MaxTokens:   100,
        Temperature: 0.7,
    }

    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        log.Fatalf("Error creating chat completion: %v", err)
    }

    if len(resp.Choices) > 0 {
        fmt.Println("\nText Response from Groq API:\n")
        fmt.Println(resp.Choices[0].Message.Content)
    }
}

func chatCompletionCreateJSON(ctx context.Context, client *groq.Client, modelID string) {
    req := groq.ChatCompletionRequest{
        Model: modelID,
        Messages: []groq.Message{
            {Role: "user", Content: "What is Golang? Respond in JSON format with keys: 'name', 'description', and 'key_features'."},
        },
        MaxTokens:   150,
        Temperature: 0.7,
        ResponseFormat: &groq.ResponseFormat{
            Type: "json_object",
        },
    }

    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        log.Fatalf("Error creating chat completion (JSON): %v", err)
    }

    if len(resp.Choices) > 0 {
        fmt.Println("\nJSON Response from Groq API:\n")
        fmt.Println(resp.Choices[0].Message.Content)
    } else {
        fmt.Println("No response received from API")
    }
}

func chatCompletionCreateStream(ctx context.Context, client *groq.Client, modelID string) {
    req := groq.ChatCompletionRequest{
        Model: modelID,
        Messages: []groq.Message{
            {Role: "user", Content: "Tell me a short story about a Tesla Optimus."},
        },
        MaxTokens:   150,
        Temperature: 0.7,
        Stream:      true,
    }

    fmt.Println("\nStreaming Response from Groq API:\n")
    chunkChan, errChan := client.CreateChatCompletionStream(ctx, req)
    for {
        select {
        case chunk, ok := <-chunkChan:
            if !ok {
                return
            }
            for _, choice := range chunk.Choices {
                fmt.Print(choice.Delta.Content)
            }
        case err, ok := <-errChan:
            if !ok {
                return
            }
            if err != nil {
                log.Fatalf("Error in stream: %v", err)
            }
            return
        }
    }
}
```

## Features

- Simple and intuitive API
- Support for chat completions
- Streaming responses
- Model information retrieval
- Environment variable configuration with godotenv support
- Configurable base URL and timeout

## Documentation

For detailed API documentation, please refer to the [godoc](https://pkg.go.dev/github.com/wansatya/groq-go) page.

## Development

### Prerequisites

- Go 1.16 or higher
- Make (for running Makefile commands)

### Building

To build the project, run:

```
make build
```

### Testing

To run the tests, use:

```
make test
```

### Linting

To lint the code, use:

```
make lint
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

For more details, see the [CONTRIBUTING](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial SDK and is not affiliated with or endorsed by Groq, Inc.