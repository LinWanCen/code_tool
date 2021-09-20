package util

import (
	"path/filepath"
	"regexp"
	"strings"
)

func FileNameHasPrefix(path string, prefixes []string) (fileName string, needParse bool) {
	fileName = filepath.Base(path)
	for _, prefix := range prefixes {
		needParse = needParse || strings.HasPrefix(fileName, prefix)
	}
	return
}

func RelativePath(rootPath string, path string) string {
	return path[len(rootPath):]
}

var mavenRegExp = regexp.MustCompile(`([^/\\]*)[/\\]src[/\\](?:main|test)[/\\]java[/\\](.*)\.java`)

func MavenPath(path string) (project string, className string, classSimpleName string) {
	submatch := mavenRegExp.FindStringSubmatch(path)
	if submatch == nil {
		base := filepath.Base(path)
		classSimpleName = base[0 : len(base)-len(filepath.Ext(base))]
		return
	}
	project = submatch[1]
	className = strings.Replace(submatch[2], "\\", ".", -1)
	classSimpleName = filepath.Ext(className)[1:]
	return
}
