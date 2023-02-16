package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func inSlice(what rune, where []rune) bool {
	for _, r := range where {
		if r == what {
			return true
		}
	}
	return false
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// reads line-delimited contents of a file into a map of strings
func readMap(path string) map[string]bool {
	file, err := os.Open(path)
	if err != nil {
		return map[string]bool{}
	}
	defer file.Close()

	avail := map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		avail[scanner.Text()] = true
	}
	return avail
}

// writes a map of strings to a line-delimited file
func writeMap(lines map[string]bool, path string) {
	file, _ := os.Create(path)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, _ := range lines {
		fmt.Fprintln(writer, key)
	}
	writer.Flush()
}
