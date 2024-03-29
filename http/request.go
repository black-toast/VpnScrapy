package http

import (
	"io"
	"net/http"
)

type RequestConfig struct {
	Method    string
	Url       string
	Transport http.RoundTripper
	Headers   map[string]string
}

func Request(requestConfig *RequestConfig) ([]byte, error) {
	client := &http.Client{
		Transport: requestConfig.Transport,
	}
	req, err := http.NewRequest(requestConfig.Method,
		requestConfig.Url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if requestConfig.Headers != nil {
		for key, value := range requestConfig.Headers {
			req.Header.Set(key, value)
		}
	}
	SetCommonHeader(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
