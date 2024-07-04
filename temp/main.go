package main

import (
	"fmt"
	"os"
	"strings"
)

func processText(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i := 0; i < len(lines); i++ {
		currentLine := strings.TrimSpace(lines[i])

		if currentLine == "" {
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				result.WriteString("\n")
			}
		} else {
			if result.Len() > 0 {
				result.WriteString(" ")
			}
			result.WriteString(currentLine)
		}
	}

	return result.String()
}

func readAndProcessFile(inputFilename, outputFilename string) error {
	content, err := os.ReadFile(inputFilename)
	if err != nil {
		return err
	}

	processedContent := processText(string(content))
	err = os.WriteFile(outputFilename, []byte(processedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	inputFilename := "plain_text.txt"
	outputFilename := "output.txt"

	err := readAndProcessFile(inputFilename, outputFilename)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

	fmt.Println("File processed successfully.")
}
