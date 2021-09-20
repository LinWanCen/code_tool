package util

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// 把日志写到控制台和文件，并在 deep.go ForDeepListFile() 中关闭文件
var Logger *log.Logger

var logFileName = flag.String("log", outPath+"log.txt", "logFileName")

var logFile *os.File
var logWriter *bufio.Writer

func init() {
	_ = os.MkdirAll(filepath.Dir(*logFileName), 0777)
	file, err := os.OpenFile(*logFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("os.OpenFile(...): %v\n", err)
		return
	}
	logFile = file
	logWriter = bufio.NewWriter(file)
	Logger = log.New(io.MultiWriter(os.Stdout, logWriter), "", log.Lshortfile)
}

// 可以改为钩子
func FlushAndCloseLogFile() {
	if logWriter != nil {
		_ = logWriter.Flush()
	}
	if logFile != nil {
		_ = logFile.Close()
	}
}
