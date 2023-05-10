package http

import (
	"VpnScrapy/crypt"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var default_domain = "https://an.tly07.com"

var login_path = "/api/E-an2.php"

func Login(email string, password string) ([]byte, error) {
	// a=JiaMi(email)&b=password&time=curtime
	encryptResult, err := crypt.AesCBCEncrypt([]byte(email))
	if err != nil {
		return nil, err
	}
	curTimeMills := time.Now().UnixMilli()
	encryptEmail := base64.StdEncoding.EncodeToString(encryptResult)
	formBody := fmt.Sprintf("a=%s&b=%s&time=%d", encryptEmail, password, curTimeMills)

	// urli := url.URL{}
	// urlproxy, _ := urli.Parse("http://127.0.0.1:8888")
	client := &http.Client{
		// Transport: &http.Transport{
		// 	Proxy: http.ProxyURL(urlproxy),
		// },
	}
	req, err := http.NewRequest("POST",
		default_domain+login_path,
		strings.NewReader(formBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("User-Agent", "okhttp/3.2.0")
	req.Header.Set("Accept", "*/*")
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
	decodeResult, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}
	decryptResult, err := crypt.AesCBCDecrypt(decodeResult)
	if err != nil {
		return nil, err
	}
	//fmt.Println("response body:", string(decrypt_result))
	return decryptResult, nil
}

var ssrSubscribeUrl = "https://fast.lycorisrecoil.org/link/sVfp6MHIJv1auOpA?list=ssa"

func SsrSubscribe() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		ssrSubscribeUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
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
	return body, err
}
