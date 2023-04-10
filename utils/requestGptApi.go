package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goChat/conf"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

func RequestRemoteGptApi(msg string) (string, error) {
	url := conf.GetConfig("other_api", "customization_api")
	if url == "" {
		url = "http://127.0.0.1"
	}
	url = fmt.Sprintf("%s/%s", url, "/postChatMsg/caomaoboy")
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("msg", msg)
	err := writer.Close()
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func MatchChatGptPythonCode(msg string) string {
	var re *regexp.Regexp
	if strings.Contains(msg, "python") {
		re = regexp.MustCompile("```python(.*)```")
	} else if strings.Contains(msg, "sql") {
		re = regexp.MustCompile("```sql(.*)```")
	} else if strings.Contains(msg, "result") {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(msg), &data)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(data["result"], "pandasSql = ") {
			re := regexp.MustCompile(`\"(.*)\"`) // 使用正则表达式来匹配双引号中间的内容
			match := re.FindStringSubmatch(msg)  // 返回匹配到的内容，其中match[0]是完整的匹配结果，match[1]是第一个括号内的匹配结果
			if len(match) > 1 {
				return match[1]
			}
			return ""
		}
		return data["result"]
	} else {
		return msg
	}
	match := re.FindStringSubmatch(msg)
	if len(match) > 0 {

		rs := match[1]
		rs = strings.ReplaceAll(rs, "\\n\\n", "\n")
		rs = strings.ReplaceAll(rs, "\\u003c", "<")
		rs = strings.ReplaceAll(rs, "\\u003e", ">")
		rs = strings.ReplaceAll(rs, "\\n", "\n")
		rs = strings.ReplaceAll(rs, "\\\"", "\"")
		return rs
	}

	return ""
}

//MatchHandlerSqlTable 处理sql生成 `table0` 无法执行带有反引号的sql
func MatchHandlerSqlTable(msg string) string {
	// 正则表达式匹配 `table` 跟一个数字
	reg := regexp.MustCompile("`table(\\d+)`")
	newStr := reg.ReplaceAllStringFunc(msg, func(s string) string {
		return s[1 : len(s)-1]
	})

	return newStr
}

//MatchHandlerSqlTableS 处理sql生成 "table0" 无法执行带有反引号的sql
func MatchHandlerSqlTableS(msg string) string {
	// 正则表达式匹配 "table" 跟一个数字
	reg := regexp.MustCompile("\"table(\\d+)\"")
	newStr := reg.ReplaceAllStringFunc(msg, func(s string) string {
		return s[1 : len(s)-1]
	})

	return newStr
}
