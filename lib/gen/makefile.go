package gen

import (
  "fmt"
  "os"
  "path/filepath"
)

// Create a new makefile which echoes "I was just created" to the console
func CreateMakeFile(path string) {

  absPath, pathError := filepath.Abs(filepath.Join(path, "Makefile"))

  if pathError != nil {
    fmt.Println("Error creating Makefile. File path does not exist", pathError)
    os.Exit(1)
  }

  makefileContent := []byte("WATCHMEGO := ./output\n\nwatch:\n\tgo build -o output .\n")

  err := os.WriteFile(absPath, makefileContent, 0644)
  if err != nil {
    fmt.Println("Error creating Makefile", err)
    os.Exit(1)
  }

  fmt.Println("Makefile created")
}
