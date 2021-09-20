package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var outPath = "tools_out/"

func init() {
	_ = os.MkdirAll(filepath.Dir(outPath), 0777)
	file, err := os.OpenFile(outPath+".gitignore", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("os.OpenFile(...): %v\n", err)
		return
	}
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString("*")
	if logWriter != nil {
		_ = writer.Flush()
	}
	if logFile != nil {
		_ = file.Close()
	}
}
