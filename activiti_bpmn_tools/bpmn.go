package activiti_bpmn_tools

import (
	"path/filepath"
	"strings"

	"code_tools/util"
)

type (
	// Activiti 工作流 .bpmn20.xml 文件节点清单
	//  <process id="" name="">
	//    <serviceTask id="" name="">
	Bpmn struct {
		Process process `xml:"process"`
	}

	process struct {
		Id          string        `xml:"id,attr"`
		Name        string        `xml:"name,attr"`
		ServiceTask []serviceTask `xml:"serviceTask"`
	}

	serviceTask struct {
		Id   string `xml:"id,attr"`
		Name string `xml:"name,attr"`
	}
)

func (b Bpmn) FeaturesName() string {
	return "Activiti 工作流 .bpmn20.xml 文件节点清单"
}

func (b Bpmn) Before(i ...interface{}) (results []interface{}) {
	return
}

func (b Bpmn) FileNames() []string {
	return []string{"activiti_bpmn_list"}
}

func (b Bpmn) Titles() [][]string {
	return [][]string{{
		"JOB_CODE",
		"id",
		"name",
		"taskId",
		"taskName",
	}}
}

func (b Bpmn) Lines(rootPath string, path string) [][][]string {
	fileName := filepath.Base(path)
	if !strings.HasSuffix(fileName, ".bpmn20.xml") {
		return nil
	}

	model := Bpmn{}
	err := util.ParseXml(path, &model)
	if err != nil {
		return nil
	}

	var csvData [][]string
	for _, task := range model.Process.ServiceTask {
		csvData = append(csvData, []string{
			"DPS_" + model.Process.Id,
			model.Process.Id,
			model.Process.Name,
			task.Id,
			task.Name,
		})
	}
	return [][][]string{csvData}
}
