package integration

import (
	"context"
	"os"
	"testing"

	"github.com/wansatya/groq-go/pkg/groq"
)

func TestCreateChatCompletion(t *testing.T) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set, skipping integration test")
	}

	client := groq.NewClient(apiKey)

	req := groq.ChatCompletionRequest{
		Model: "mixtral-8x7b-32768",
		Messages: []groq.Message{
			{Role: "user", Content: "Say 'Hello, World!'"},
		},
		MaxTokens:   10,
		Temperature: 0,
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("Error creating chat completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices returned in response")
	}

	content := resp.Choices[0].Message.Content
	if content != "Hello, World!" {
		t.Errorf("Unexpected response content: %s", content)
	}
}