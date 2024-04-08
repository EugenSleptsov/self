package main

import (
	"fmt"
	"os"
)

func main() {
	content := ``

	file, err := os.OpenFile("main.go", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	file.WriteString(content)
}
