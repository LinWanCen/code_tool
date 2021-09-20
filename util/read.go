package util

import (
	"io/ioutil"
	"os"
)

// 读取文件
func ReadByte(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		Logger.Printf("os.Open(...): %v\n", err)
		return nil
	}
	defer func() {
		if err := file.Close(); err != nil {
			Logger.Printf("os.Close(...): %v\n", err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		Logger.Printf("ioutil.ReadAll(...): %v\n", err)
		return nil
	}
	return data
}
