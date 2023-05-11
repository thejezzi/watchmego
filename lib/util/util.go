package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskYesNoQuestion(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt + " (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "y", "Y":
			return true
		case "n", "N":
			return false
		default:
			fmt.Println("Please enter y or n")
		}
	}
}

func ReadMakeFileVar(identifier string) string {
	file, err := os.Open("Makefile")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, identifier) {
			return strings.Split(line, "=")[1]
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return ""
}

func ParseArguments(args string) []string {
	var parsedArgs []string
	var currentArg string
	var inQuotes bool

	for _, char := range args {
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}

		if char == ' ' && !inQuotes {
			parsedArgs = append(parsedArgs, currentArg)
			currentArg = ""
			continue
		}

		currentArg += string(char)
	}

	parsedArgs = append(parsedArgs, currentArg)

	return parsedArgs
}
