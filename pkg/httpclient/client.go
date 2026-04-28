package httpclient

import (
	"net/http"
	"time"
)

func GetExternalAPIClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        50, // 外网并发量通常低于内网，适当减小
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second, // 外网 API 响应慢，超时给足 60 秒
	}
}
