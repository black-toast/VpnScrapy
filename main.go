package main

import (
	"VpnScrapy/crawl/novel"
	"VpnScrapy/http"
	"fmt"
)

func main() {
	//server.ConfigServer()
	novelContent, err := http.Request(&http.RequestConfig{
		Method:    "GET",
		Url:       novel.WithUrl().BuildNovelUrl(),
		Transport: http.V2rayProxy(),
	})
	if err != nil {
		fmt.Printf("request %s fail\n", novel.WithUrl().BuildNovelUrl())
		return
	}
	novel.Scrapy(string(novelContent))
}
