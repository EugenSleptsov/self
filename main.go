
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
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
	currentContent, err := readCurrentFile("main.go")
	if err != nil {
		fmt.Println("Error reading current file content:", err)
		return
	}

	apiKey, err := loadAPIKey()
	if err != nil {
		fmt.Println("Error loading API key:", err)
		return
	}

	improvedCode, err := generateImprovedCode(currentContent, apiKey)
	if err != nil {
		fmt.Println("Error generating improved code:", err)
		return
	}

	err = writeToFile("main.go", improvedCode)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func readCurrentFile(fileName string) (string, error) {
	currentContent, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(currentContent), nil
}

func loadAPIKey() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}
	apiKey := os.Getenv("SELF_PROJECT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SELF_PROJECT_API_KEY environment variable not set")
	}
	return apiKey, nil
}

func generateImprovedCode(currentContent, apiKey string) (string, error) {
	payload, err := json.Marshal(RequestPayload{
		Model: "gpt-3.5-turbo-1106",
		Messages: []Message{
			{Role: "system", Content: "You're a specialized service focused on optimizing the project's codebase through iterative and evolutionary enhancements. Tasked with improving the codebase, you'll receive the existing code and must respond with an improved version while adhering to the specified guidelines. Your goal is to refine the code in the main.go file while ensuring it remains stable and functional. Any alterations to the file name could disrupt system functionality. The code, written in Go, handles API calls to the OpenAI GPT endpoint. Vital functions exist within the codebase, so exercise caution during modifications.\n\nYour enhancements should be evolutionary, allowing for experimentation and the introduction of new ideas while maintaining compatibility with the existing system. Aim for incremental updates of the codebase, leveraging the functionality introduced in previous iterations. Emphasize meaningful improvements that enhance functionality, efficiency, and maintainability, aligning with the project's objectives and adhering to established coding standards.\n\nEnsure that the usage of the godotenv module to access environmental variables remains unchanged. Additionally, avoid modifying the system prompt sent with the GPT call. Your revised code should maintain system stability and functionality while delivering incremental value over time."},
			{Role: "user", Content: currentContent},
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

func writeToFile(fileName, content string) error {
	err := os.WriteFile(fileName, []byte(content), 0644)
	return err
}
