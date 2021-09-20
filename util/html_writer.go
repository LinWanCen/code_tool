package util

import (
	"bufio"
	"html/template"
	"os"
	"path/filepath"
)

// 分隔符文件写入
type HtmlWriter struct {
	*os.File
	*bufio.Writer
}

func NewHtmlWriter(fileName string, htmlTitle string, title []string) *HtmlWriter {
	c := new(HtmlWriter)

	_ = os.MkdirAll(filepath.Dir(fileName), 0777)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		Logger.Printf("os.OpenFile(...): %v\n", err)
		return c
	}
	writer := bufio.NewWriter(file)
	// 写入 HTML 内容前缀
	_, err = writer.WriteString(pre1)
	_, err = writer.WriteString(htmlTitle)
	_, err = writer.WriteString(pre2)
	if err != nil {
		Logger.Printf("writer.WriteString(...): %v\n", err)
		if err := file.Close(); err != nil {
			Logger.Printf("file.Close(): %v\n", err)
		}
		return c
	}
	// 写入标题
	err = titleTemp.Execute(writer, title)
	if err != nil {
		Logger.Printf("titleTemp.Execute(...): %v\n", err)
		if err := file.Close(); err != nil {
			Logger.Printf("file.Close(): %v\n", err)
		}
		return c
	}
	c.File = file
	c.Writer = writer
	return c
}

func (c *HtmlWriter) Write(data [][]string) {
	if data == nil {
		return
	}
	if c.Writer == nil {
		return
	}
	err := lineTemp.Execute(c.Writer, data)
	if err != nil {
		Logger.Printf("lineTemp.Execute(...): %v\n", err)
	}
}

func (c *HtmlWriter) Close() {
	if c.Writer != nil {
		// 写入 HTML 内容后缀
		_, err := c.Writer.WriteString(suffix)
		if err != nil {
			Logger.Printf("c.Writer.WriteString(...): %v\n", err)
		}
		err = c.Writer.Flush()
		if err != nil {
			Logger.Printf("c.Writer.Flush(): %v\n", err)
		}
	}
	if c.File != nil {
		if err := c.File.Close(); err != nil {
			Logger.Printf("c.File.Close(): %v\n", err)
		}
	}
}

var titleTemp = template.Must(template.New("titleTemp").Parse(
	// language="GoTemplate"
	`  <tr>{{range $index, $value := .}}<th>{{$value}}</th>{{end}}</tr>
`))

var lineTemp = template.Must(template.New("lineTemp").Parse(
	// language="GoTemplate"
	`{{range $i, $tr := .}}
  <tr>{{range $j, $td := $tr}}<td>{{$td}}</td>{{end}}</tr>{{end}}`))

const pre1 =
// language="html"
`<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="UTF-8">
  <title>`

const pre2 =
// language="html"
`</title>
</head>

<style>
table {
  border-collapse: collapse;
}

th, td {
  border: 1px solid lightgray;
  padding-left: 5px;
  padding-right: 5px;
}
</style>

<body>
<table>
`

const suffix =
// language="html"
`
</table>
</body>

</html>
`
