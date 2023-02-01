package main

import (
	"bufio"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"os"
)

const API_KEY = "API_KEY"

func main() {
	key := os.Getenv(API_KEY)
	if key == ""{
		log.Fatalf("failed to find %v", API_KEY)
	}
	c := gogpt.NewClient(key)
	ctx := context.Background()

	for {
		fmt.Printf("Human: \n")
		f := bufio.NewReader(os.Stdin)
		input, _ := f.ReadString('\n')
		req := gogpt.CompletionRequest{
			Model:           gogpt.GPT3TextDavinci003,
			MaxTokens:       1500,
			Temperature:     0.9,
			TopP:            1,
			BestOf:          1,
			PresencePenalty: 0.6,
			Prompt:          input,
		}

		resp, err := c.CreateCompletion(ctx, req)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			return
		}
		fmt.Printf("AI: ")
		for _, r := range resp.Choices {
			fmt.Printf(r.Text)
		}
		fmt.Println("")
	}
}
