package util

import (
	"flag"
	"os"
	"path/filepath"
)

type DeepHandle interface {
	BeforeEveryPath(rootPath string)
	EveryPath(rootPath, path string)
	AfterEveryPath(rootPath string)
}

// 读取第一个参数或拖动到exe的文件夹，遍历所有子文件
func ForDeepListFileFlag(handles []DeepHandle) {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		path = "."
	}
	ForDeepListFile(path, handles)
}

func ForDeepListFile(rootPath string, handles []DeepHandle) {
	for _, writer := range handles {
		writer.BeforeEveryPath(rootPath)
	}
	err := filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		for _, writer := range handles {
			writer.EveryPath(rootPath, path)
		}
		return nil
	})
	if err != nil {
		Logger.Printf("filepath.Walk(...): %v\n", err)
	}
	for _, writer := range handles {
		writer.AfterEveryPath(rootPath)
	}
	FlushAndCloseLogFile()
}
