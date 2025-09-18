package util

import (
	"fmt"
	"io"
	"net/http"
	"onij/util/logs"

)

func Get(url string, param map[string]string) ([]byte, error) {
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

	client := &http.Client{}
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