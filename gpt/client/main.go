package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var url = "http://localhost:8080/"

type InputParam struct {
	InputParam string `json:"inputParam"`
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	for {

		fmt.Printf("Human: \n")
		f := bufio.NewReader(os.Stdin)
		str, _ := f.ReadString('\n')

		input := InputParam{InputParam: str}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			log.Fatalf("failed to convert json to struct, %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(inputJSON))
		if err != nil {
			log.Fatalf("failed to get response, %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to get response body, %v", err)
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Printf("error decoding response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("response: %q", body)
		}
		fmt.Printf("AI:%s\n",response.Message)

	}
}
