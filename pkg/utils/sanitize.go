package utils

import "strings"

func TrimString(s string) string {
	return strings.TrimRight(s, "\r\n")
}
