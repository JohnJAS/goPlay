package main

import (
	"context"
	"encoding/json"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"net/http"
	"os"
)

type InputParam struct {
	InputParam string `json:"inputParam"`
}

type Response struct {
	Message string `json:"message"`
}

const API_KEY = "API_KEY"

func main() {
	key := os.Getenv(API_KEY)
	if key == "" {
		log.Fatalf("failed to find %v", API_KEY)
	}
	c := gogpt.NewClient(key)
	ctx := context.Background()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var input InputParam
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var result string
			result, err = callGPT(c, &ctx, input.InputParam)
			if err != nil {
				log.Printf("Error: %v, Result: %v", err, result)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response := Response{Message: fmt.Sprintf("%s", result)}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func callGPT(c *gogpt.Client, ctx *context.Context, input string) (string, error) {
	req := gogpt.CompletionRequest{
		Model:           gogpt.GPT3TextDavinci003,
		MaxTokens:       1500,
		Temperature:     0.9,
		TopP:            1,
		BestOf:          1,
		PresencePenalty: 0.6,
		Prompt:          input,
	}

	resp, err := c.CreateCompletion(*ctx, req)
	if err != nil {
		return "", err
	}

	var result string
	for _, r := range resp.Choices {
		result = result + fmt.Sprintf(r.Text)
	}

	return result, nil
}
