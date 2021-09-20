package util

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// 分隔符文件写入
type Writer struct {
	*os.File
	*csv.Writer
}

func NewWriter(fileName string, comma rune, title []string) *Writer {
	c := new(Writer)

	_ = os.MkdirAll(filepath.Dir(fileName), 0777)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		Logger.Printf("os.OpenFile(...): %v\n", err)
		return c
	}
	writer := csv.NewWriter(file)
	writer.Comma = comma
	err = writer.Write(title)
	if err != nil {
		Logger.Printf("writer.Write(...): %v\n", err)
		if err := file.Close(); err != nil {
			Logger.Printf("file.Close(...): %v\n", err)
		}
		return c
	}
	c.File = file
	c.Writer = writer
	return c
}

func (c *Writer) Write(data [][]string) {
	if data == nil {
		return
	}
	if c.Writer == nil {
		return
	}
	if err := c.Writer.WriteAll(data); err != nil {
		Logger.Printf("c.Writer.WriteAll(...): %v\n", err)
	}
}

func (c *Writer) Close() {
	if c.Writer != nil {
		c.Writer.Flush()
	}
	if c.File != nil {
		if err := c.File.Close(); err != nil {
			Logger.Printf("c.File.lose(...): %v\n", err)
		}
	}
}
