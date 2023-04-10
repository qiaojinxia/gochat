package utils

import (
	"os"
	"strings"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func FilterTrim(content string) string {
	content = strings.TrimSuffix(content, "\r\n")
	content = strings.TrimSuffix(content, "\r")
	content = strings.TrimSuffix(content, "\n")
	return content
}
