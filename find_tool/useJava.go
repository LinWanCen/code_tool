package mybatis_tool

import (
	"path/filepath"
	"regexp"

	"code_tools/util"
)

type (
	Java struct {
	}
)

func (j Java) FeaturesName() string {
	return "Java 分析"
}

func (j Java) Before(i ...interface{}) (results []interface{}) {
	return
}

func (j Java) FileNames() []string {
	return []string{"mybatis_java_list"}
}

func (j Java) Titles() [][]string {
	return [][]string{
		{"id", "project", "classSimpleName", "className"},
	}
}

var idRegexp = regexp.MustCompile(`"([a-zA-Z]+\.[a-zA-Z]+)"`)

func (j Java) Lines(rootPath string, path string) [][][]string {
	if filepath.Ext(path) != ".java" {
		return nil
	}

	m := idRegexp.FindAllSubmatch(util.ReadByte(path), -1)
	if len(m) == 0 {
		return nil
	}

	project, className, classSimpleName := util.MavenPath(path)

	result := [][][]string{{}}
	for _, mm := range m {
		result[0] = append(result[0], []string{
			string(mm[1]),
			project,
			classSimpleName,
			className,
		})
	}
	return result
}
