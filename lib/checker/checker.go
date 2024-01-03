package checker

import (
	"bufio"
	"os"
	"strings"

	"github.com/thejezzi/watchmego/lib/logger"
)

// check if makefile exists for specified directory
func CheckMakefile(dir string) bool {
	_, err := os.Stat(dir + "/Makefile")
	return err == nil
}

// Read Makefile and check if there is an entry for "watch"
func CheckMakefileWatch(dir string) bool {
	file, err := os.Open(dir + "/Makefile")
	if err != nil {
		logger.Error("Error opening Makefile")
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "watch-build:") {
			return true
		}
	}
	return false
}
