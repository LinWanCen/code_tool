package main

import (
	"code_tools/mybatis_tool"
	"code_tools/util"
)

func main() {
	util.ForFlag([]util.ToLine{
		new(mybatis_tool.Mapper),
		new(mybatis_tool.Java),
	})
}
