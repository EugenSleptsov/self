
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ApiResponse struct {
	Choices []Choice `json:"choices"`
}

func main() {
	// Read current main.go file content
	currentContent, err := ioutil.ReadFile("main.go")
	if err != nil {
		fmt.Println("Error reading current file content:", err)
		return
	}

	// Make API call to OpenAI GPT endpoint to improve the code
	improvedCode, err := improveCode(string(currentContent))
	if err != nil {
		fmt.Println("Error improving code:", err)
		return
	}

	// Write improved code to file
	err = ioutil.WriteFile("main.go", []byte(improvedCode), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func improveCode(currentContent string) (string, error) {
	// Retrieve API key from environment variable
	apiKey := os.Getenv("SELF_PROJECT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SELF_PROJECT_API_KEY environment variable not set")
	}

	// Create payload
	payload := RequestPayload{
		Model: "gpt-3.5-turbo-1106",
		Messages: []Message{
			{Role: "system", Content: "You are a service that improves code of the project. I will send you a code and you need to answer with improved code. Answer with improved code only. Your code must be in one file. If this file would fail to run, then it will break everything, so try your best to not break anything. The code should be in Go. The code is the current project that handles the API call to OpenAI GPT endpoint. You should never change a Model that is used in payload and endpoint url"},
			{Role: "user", Content: currentContent},
		},
	}

	// Make POST request to OpenAI GPT endpoint
	improvedContent, err := makePostRequest(payload, apiKey)
	if err != nil {
		return "", err
	}

	// Extract improved code from API response
	improvedCode := extractImprovedCode(improvedContent)

	return improvedCode, nil
}

func extractImprovedCode(content string) string {
	// Use regex to extract code between ```go tags
	re := regexp.MustCompile("```go([\\s\\S]*)```")
	match := re.FindStringSubmatch(content)
	if len(match) >= 2 {
		return match[1]
	}
	return content
}

func makePostRequest(payload RequestPayload, apiKey string) (string, error) {
	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Make HTTP POST request to OpenAI GPT endpoint
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client
	client := http.Client{}

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
