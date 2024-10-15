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