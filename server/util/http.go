package util

import (
	"fmt"
	"io"
	"net/http"
	"onij/util/logs"
	"path"
	"strings"
	"time"
)

func GET(url string, param map[string]string) ([]byte, error) {
	for k, v := range param {
		if k == "" || v == "" {
			continue
		}
		url += fmt.Sprintf("&%s=%s", k, v)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetUrlResourceName(url string) string {
	// 去除查询参数和片段标识
	cleanURL := strings.Split(url, "?")[0]
	cleanURL = strings.Split(cleanURL, "#")[0]

	// 获取路径的最后一部分
	resource := path.Base(cleanURL)

	// 如果结果是空或者斜杠，返回空字符串
	if resource == "" || resource == "/" || resource == "." {
		return ""
	}

	return resource
}
