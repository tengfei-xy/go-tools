package tools

import (
	"fmt"
	"io"
	"net/http"
)

func sendGetRequest(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("url:%s StatusCode:%d", link, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
