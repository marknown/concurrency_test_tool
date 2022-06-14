package utils

import (
	"concurrency_test_tool/libs/parser"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/marknown/utils"
)

func BuildRequest(source string) (*sling.Sling, error) {
	parseResult := parser.NewCurlParser(source).Parse()

	if parseResult.URL == "" {
		return nil, errors.New("没有解析到URL，请核实是否为完整的CURL命令")
	}

	s := sling.New()

	// 设置请求头
	for k, v := range parseResult.Headers {
		s.Set(k, v)
	}

	// 设置请求方式
	switch parseResult.Method {
	case "GET":
		s.Get(parseResult.URL)
	case "POST":
		s.Post(parseResult.URL)
	case "HEAD":
		s.Head(parseResult.URL)
	case "PUT":
		s.Put(parseResult.URL)
	case "PATCH":
		s.Patch(parseResult.URL)
	case "DELETE":
		s.Delete(parseResult.URL)
	case "OPTIONS":
		s.Options(parseResult.URL)
	case "TRACE":
		s.Trace(parseResult.URL)
	case "CONNECT":
		s.Connect(parseResult.URL)
	default:
		s.Get(parseResult.URL)
	}

	if v, ok := parseResult.Args["--data-raw"]; ok {
		s.Body(strings.NewReader(v))
	}

	if v, ok := parseResult.Args["--data"]; ok {
		s.Body(strings.NewReader(v))
	}

	if v, ok := parseResult.Args["-d"]; ok {
		s.Body(strings.NewReader(v))
	}

	return s, nil
}

func DoRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// when err is nil, resp contains a non-nil resp.Body which must be closed
	defer resp.Body.Close()

	// The default HTTP client's Transport may not
	// reuse HTTP/1.x "keep-alive" TCP connections if the Body is
	// not read to completion and closed.
	// See: https://golang.org/pkg/net/http/#Response
	defer io.Copy(ioutil.Discard, resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// ParseTimeFromString 把字符串日期解析成 time.Time 类型
func ParseTimeFromString(t string) (time.Time, error) {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.ParseInLocation("2006-01-02 15:04:05", t, cstZone)
}

// FormatTimeStringCNSecond format cn 2006-01-02 15:04:05.999
func FormatTimeStringCNSecond(t time.Time) string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	return t.In(cstZone).Format("2006-01-02 15:04:05")
}

// FormatTimeStringCNMillisecond format cn 2006-01-02 15:04:05.999
func FormatTimeStringCNMillisecond(t time.Time) string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	return t.In(cstZone).Format("2006-01-02 15:04:05.999")
}

// UnixTimstampMillisecond 返回毫秒级时间戳
func UnixTimstampMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetCurrentDirectory 获取可执行文件的当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}
func LogInfo(info string) {
	path := GetCurrentDirectory() + "/concurrency_log.txt"
	utils.FileAppend(path, info)
}
