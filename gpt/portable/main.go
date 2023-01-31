package main

import (
	"bufio"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"os"
)

const API_KEY = "API_KEY"

func main() {
	c := gogpt.NewClient(os.Getenv(API_KEY))
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
			return
		}
		fmt.Printf("AI: ")
		for _, r := range resp.Choices {
			fmt.Printf(r.Text)
		}
		fmt.Println("")
	}
}
