package main

func inSlice(what rune, where []rune) bool {
	for _, r := range where {
		if r == what {
			return true
		}
	}
	return false
}

