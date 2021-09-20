package util

import "strings"

func FileLink(path string) string {
	return "file:///" + strings.ReplaceAll(path, `\`, `/`)
}

func LineLink(path, link string) string {
	return "file:///" + strings.ReplaceAll(path, `\`, `/`) + ":" + link
}
