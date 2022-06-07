package util

import "strings"

func Empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
