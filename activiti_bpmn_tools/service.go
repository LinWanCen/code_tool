package activiti_bpmn_tools

import (
	"path/filepath"

	"code_tools/util"
)

// Activiti 节点对应类和方法清单
//  <request name="" service="" method="">
type Service struct {
	Name    string `xml:"name,attr"`
	Service string `xml:"service,attr"`
	Method  string `xml:"method,attr"`
}

func (s Service) FeaturesName() string {
	return "Activiti 节点对应类和方法清单"
}

func (s Service) Before(i ...interface{}) (results []interface{}) {
	return
}

func (s Service) FileNames() []string {
	return []string{"activiti_service_list"}
}

func (s Service) Titles() [][]string {
	return [][]string{{
		"JOB_CODE",
		"name",
		"service",
		"method",
	}}
}

func (s Service) Lines(rootPath string, path string) [][][]string {
	if filepath.Ext(path) != ".xml" {
		return nil
	}

	fileName, done := util.FileNameHasPrefix(path, []string{
		"SS",
		"TR",
	})
	if !done {
		return nil
	}

	// 删除拓展名
	fileName = fileName[0 : len(fileName)-4]

	model := Service{}
	err := util.ParseXml(path, &model)
	if err != nil {
		return nil
	}

	return [][][]string{{{
		fileName,
		model.Name,
		model.Service,
		model.Method,
	}}}
}
