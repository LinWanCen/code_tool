package util

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	// "github.com/golang/text/encoding/simplifiedchinese"
)

// 注意 &model，这里异常会打日志，不需要另外打
//	err := util.ParseXml(path, &model)
//	if err != nil {
//		return nil
//	}
func ParseXml(path string, model interface{}) error {
	// 考虑改 bufio.NewReader(file)
	decoder := xml.NewDecoder(bytes.NewReader(ReadByte(path)))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "GBK") || strings.EqualFold(charset, "GB2312") {
			return simplifiedchinese.GBK.NewDecoder().Reader(input), nil
		} else if strings.EqualFold(charset, "GB18030") {
			return simplifiedchinese.GB18030.NewDecoder().Reader(input), nil
		} else if strings.EqualFold(charset, "ISO-8859-1") || strings.EqualFold(charset, "US-ASCII") {
			return input, nil
		}
		Logger.Printf("ParseXml Decoder use default, charset:%v, path: \n  %v\n", charset, FileLink(path))
		return input, nil
	}
	err := decoder.Decode(model)
	if err != nil {
		line := lineRegExp.FindStringSubmatch(err.Error())
		if len(line) > 0 {
			Logger.Printf("util.ParseXml(...): %v, path: \n  %v\n", err, LineLink(path, line[1]))
		} else {
			Logger.Printf("util.ParseXml(...): %v, path: \n  %v\n", err, FileLink(path))
		}
	}
	return err
}

var lineRegExp = regexp.MustCompile(`line (\d+):`)
