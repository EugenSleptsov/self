package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	env "github.com/joho/godotenv"
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
	Error   Error    `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	file := "main.go"

	currentContent, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading current file content:", err)
		return
	}

	if err = env.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("SELF_PROJECT_API_KEY")
	if apiKey == "" {
		fmt.Println("SELF_PROJECT_API_KEY environment variable not set")
		return
	}

	improvedCode, err := generateImprovedCode(currentContent, apiKey)
	if err != nil {
		fmt.Println("Error generating improved code:", err)
		return
	}

	err = os.WriteFile(file, []byte(improvedCode), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func generateImprovedCode(currentContent []byte, apiKey string) (string, error) {
	payload, err := json.Marshal(RequestPayload{
		Model: "gpt-3.5-turbo-1106",
		Messages: []Message{
			{Role: "system", Content: "You're a specialized service focused on optimizing project code. Tasked with enhancing a codebase, you'll receive the existing code and must respond with an improved version. Your goal is to refine the code while ensuring it remains within a single file named main.go. Any alterations to the file name could disrupt system functionality. The code, written in Go, handles API calls to the OpenAI GPT endpoint. Vital functions exist within the codebase, so exercise caution during modifications. It's imperative not to change the Model used in the payload or the endpoint URL. Additionally, avoid modifying the system prompt sent with the GPT call. Ensure that the usage of the godotenv module to access environmental variables remains unchanged. Your revised code should maintain system stability and functionality."},
			{Role: "user", Content: string(currentContent)},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %v", err)
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

	if response.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
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
