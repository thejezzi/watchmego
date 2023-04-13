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
