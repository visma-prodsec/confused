package main

import "strings"

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
