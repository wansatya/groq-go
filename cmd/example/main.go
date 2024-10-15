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

    // Add system prompts
    client.AddSystemPrompt("You are a helpful assistant. Always be polite and concise.")
    client.AddSystemPrompt("Provide examples when explaining concepts.")

    // Example with default text response
    reqText := groq.ChatCompletionRequest{
        Model: "mixtral-8x7b-32768",
        Messages: []groq.Message{
            {Role: "user", Content: "What is Golang?"},
        },
        MaxTokens:   100,
        Temperature: 0.7,
    }

    // Example with JSON response
    reqJSON := groq.ChatCompletionRequest{
        Model: "mixtral-8x7b-32768",
        Messages: []groq.Message{
            {Role: "user", Content: "What is Golang? Respond in JSON format."},
        },
        MaxTokens:   512,
        Temperature: 0.25,
        ResponseFormat: &groq.ResponseFormat{
            Type: "json_object",
        },
    }

    ctx := context.Background()

    // Make request with text response
    respText, err := client.CreateChatCompletion(ctx, reqText)
    if err != nil {
        log.Fatalf("Error creating chat completion (text): %v", err)
    }

    if len(respText.Choices) > 0 {
        fmt.Println("Text Response from Groq API:\n")
        fmt.Println(respText.Choices[0].Message.Content)
    }

    // Make request with JSON response
    respJSON, err := client.CreateChatCompletion(ctx, reqJSON)
    if err != nil {
        log.Fatalf("Error creating chat completion (JSON): %v", err)
    }

    if len(respJSON.Choices) > 0 {
        fmt.Println("\nJSON Response from Groq API:\n")
        fmt.Println(respJSON.Choices[0].Message.Content)
    }
}