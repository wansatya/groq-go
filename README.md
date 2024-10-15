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

```
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

    "github.com/wansatya/groq-go/pkg/groq"
    "github.com/joho/godotenv"
)

func init() {
  // Load the .env file in the current directory
  err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }
}

func main() {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        log.Fatal("GROQ_API_KEY not found in environment variables")
    }

    client := groq.NewClient(apiKey)

    req := groq.ChatCompletionRequest{
        Model: "mixtral-8x7b-32768",
        Messages: []groq.Message{
            {Role: "user", Content: "What is Golang?"},
        },
        MaxTokens:   100,
        Temperature: 0.7,
    }

    ctx := context.Background()
    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        log.Fatalf("Error creating chat completion: %v", err)
    }

    if len(resp.Choices) > 0 {
        fmt.Println("Response from Groq API:")
        fmt.Println(resp.Choices[0].Message.Content)
    } else {
        fmt.Println("No response received from API")
    }
}
```

## Features

- Simple and intuitive API
- Support for chat completions
- Configurable base URL and timeout
- Context support for cancellation and timeouts

## Documentation

For detailed API documentation, please refer to the [API.md](docs/API.md) file in the docs directory.

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

For more details, see the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial SDK and is not affiliated with or endorsed by Groq, Inc.