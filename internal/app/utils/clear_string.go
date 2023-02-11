package utils

import "strings"

func ClearString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, "\n")
	s = strings.ToLower(s)

	return s
}
