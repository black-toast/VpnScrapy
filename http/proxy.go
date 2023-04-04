package http

import (
	"net/http"
	"net/url"
)

func NilProxy() http.RoundTripper {
	return nil
}

func FiddlerProxy() http.RoundTripper {
	urlProxy, _ := new(url.URL).Parse("http://127.0.0.1:8888")
	return &http.Transport{
		Proxy: http.ProxyURL(urlProxy),
	}
}

func V2rayProxy() http.RoundTripper {
	urlProxy, _ := new(url.URL).Parse("http://127.0.0.1:10809")
	return &http.Transport{
		Proxy: http.ProxyURL(urlProxy),
	}
}
