package util

import (
	"flag"
	"fmt"
	"html/template"
)

// 第一层数组是处理的一个文件生成多个文件
// 第二层数组是生成的一个文件写入多行
// 第三层数组是一行里的每个字段
type ToLine interface {
	FeaturesName() string
	// 用于结构体中 func 成员赋值供子类重写等情况
	Before(...interface{}) (results []interface{})
	FileNames() []string
	Titles() [][]string
	Lines(rootPath, path string) [][][]string
}

type LineWriter struct {
	ToLine
	csvWriter  []*Writer
	htmlWriter []*HtmlWriter
}

func (x LineWriter) BeforeEveryPath(rootPath string) {
	x.Before()
}

func (x LineWriter) EveryPath(rootPath string, path string) {
	lines := x.Lines(rootPath, path)
	if lines == nil || len(lines) == 0 {
		return
	}
	for i := range x.csvWriter {
		if lines[i] == nil || len(lines[i]) == 0 {
			continue
		}
		x.csvWriter[i].Write(lines[i])
		x.htmlWriter[i].Write(lines[i])
	}
}

func (x LineWriter) AfterEveryPath(rootPath string) {
	for i := range x.csvWriter {
		x.csvWriter[i].Close()
		x.htmlWriter[i].Close()
	}
}

// func For(path string, comma rune, toLines []ToLine) {
// 	ForDeepListFile(path, toLineDeepHandle(toLines, comma))
// }

var comma = flag.String("c", "\t", "comma")

func ForFlag(toLines []ToLine) {
	// 先在 comma 指针写入值，ForDeepListFileFlag 里再 Parse 不影响
	flag.Parse()
	if comma == nil || *comma == `\t` {
		*comma = "\t"
	}
	ForDeepListFileFlag(toLineDeepHandle(toLines, rune((*comma)[0])))
}

func toLineDeepHandle(xmlToLines []ToLine, comma rune) []DeepHandle {
	var writers []DeepHandle
	// 遍历工具
	indexWriter := NewHtmlWriter(outPath+"index.html", "tools_out", []string{"FeaturesName", "htmlName", "csvName"})
	for _, xmlToLine := range xmlToLines {
		writer := LineWriter{}
		writer.ToLine = xmlToLine
		// 遍历工具的多个文件
		for i, fileName := range xmlToLine.FileNames() {
			csvName := fileName + ".txt"
			writer.csvWriter = append(writer.csvWriter, NewWriter(outPath+csvName, comma, xmlToLine.Titles()[i]))
			htmlName := fileName + ".html"
			writer.htmlWriter = append(writer.htmlWriter, NewHtmlWriter(outPath+htmlName, fileName, xmlToLine.Titles()[i]))
			err := indexLineTemp.Execute(indexWriter.Writer, map[string]string{
				"FeaturesName": xmlToLine.FeaturesName(),
				"htmlName":     htmlName,
				"csvName":      csvName,
			})
			if err != nil {
				msg := fmt.Sprintf("indexLineTemp.Execute(...): %v\n", err)
				Logger.Printf(msg)
				panic(msg)
			}
		}
		writers = append(writers, writer)
	}
	// 这个链接在表内却没有格子，会跑到上面
	_, err := indexWriter.Writer.WriteString(`<a href="log.txt">log.txt</a>`)
	if err != nil {
		msg := fmt.Sprintf("indexWriter.Writer.WriteString(...): %v\n", err)
		Logger.Printf(msg)
		panic(msg)
	}
	indexWriter.Close()
	return writers
}

var indexLineTemp = template.Must(template.New("").Parse(
	// language="html"
	`  <tr>
    <td>{{.FeaturesName}}</td>
    <td><a href="{{.htmlName}}">{{.htmlName}}</a></td>
    <td><a href="{{.csvName}}">{{.csvName}}</a></td>
  </tr>`))
