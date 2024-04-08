
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	currentContent, err := os.ReadFile("main.go")
	if err != nil {
		fmt.Println("Error reading current file content:", err)
		return
	}

	improvedCode, err := generateImprovedCode(currentContent)
	if err != nil {
		fmt.Println("Error generating improved code:", err)
		return
	}

	err = os.WriteFile("main.go", []byte(improvedCode), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func generateImprovedCode(currentContent []byte) (string, error) {
	payload, err := json.Marshal(RequestPayload{
		Model: "gpt-3.5-turbo-1106",
		Messages: []Message{
			{Role: "system", Content: "You are a service that improves code of the project. I will send you a code and you need to answer with improved code. Answer with improved code only. Your code must be in one file and you should not change name of the file main.go. If this file would fail to run, then it will break everything, so try your best to not break anything. The code should be in Go. The code is the current project that handles the API call to OpenAI GPT endpoint. You should never change a Model that is used in payload and endpoint url. There are vital functions in the code so be careful. You should not change system prompt that is sent to with GPT Call"},
			{Role: "user", Content: string(currentContent)},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %v", err)
	}

	apiKey := os.Getenv("SELF_PROJECT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SELF_PROJECT_API_KEY environment variable not set")
	}

	apiResponse, err := makePostRequest(payload, apiKey)
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}

	var response ApiResponse
	err = json.Unmarshal([]byte(apiResponse), &response)
	if err != nil {
		return "", fmt.Errorf("error parsing API response: %v", err)
	}

	improvedCode := extractImprovedCode(response.Choices[0].Message.Content)
	return improvedCode, nil
}

func makePostRequest(payload []byte, apiKey string) (string, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func extractImprovedCode(content string) string {
	re := regexp.MustCompile("```go([\\s\\S]*)```")
	match := re.FindStringSubmatch(content)
	if len(match) >= 2 {
		return match[1]
	}
	return content
}
