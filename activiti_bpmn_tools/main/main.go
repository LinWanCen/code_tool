package main

import (
	"code_tools/activiti_bpmn_tools"
	"code_tools/util"
)

func main() {
	util.ForFlag([]util.ToLine{
		new(activiti_bpmn_tools.Bpmn),
		new(activiti_bpmn_tools.Service),
	})
}
