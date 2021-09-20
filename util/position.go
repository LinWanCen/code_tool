package util

import (
	"strings"
)

// Position
// 在指定位置下面标记^
// 从1开始
func Position(str string, index int) string {
	if index <= 0 || index > len(str) {
		return str
	}
	prefixIndex := strings.LastIndex(str[:index], "\n") + 1
	suffixIndex := strings.Index(str[index:], "\n")
	spaceCount := index - prefixIndex - 1
	if spaceCount < 0 {
		spaceCount = prefixIndex - strings.LastIndex(str[:prefixIndex-1], "\n") - 1
	}
	if suffixIndex == -1 {
		return str + "\n" + strings.Repeat(" ", spaceCount) + "^"
	}
	splitIndex := index + suffixIndex + 1
	return str[:splitIndex] + strings.Repeat(" ", spaceCount) + "^\n" + str[splitIndex:]
}
