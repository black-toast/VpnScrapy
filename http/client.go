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

// var default_domain = "https://test.com"
var login_path = "/api/E-an2.php"

func Login(email string, password string) ([]byte, error) {
	// a=JiaMi(email)&b=password&time=curtime
	encrypt_result, err := crypt.AesCBCEncrypt([]byte(email))
	if err != nil {
		return nil, err
	}
	cur_timemill := time.Now().UnixMilli()
	encrypt_email := base64.StdEncoding.EncodeToString(encrypt_result)
	formbody := fmt.Sprintf("a=%s&b=%s&time=%d", encrypt_email, password, cur_timemill)
	fmt.Println("formbody=", formbody)

	// urli := url.URL{}
	// urlproxy, _ := urli.Parse("http://127.0.0.1:8888")
	client := &http.Client{
		// Transport: &http.Transport{
		// 	Proxy: http.ProxyURL(urlproxy),
		// },
	}
	req, err := http.NewRequest("POST",
		default_domain+login_path,
		strings.NewReader(formbody),
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	decode_result, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}
	decrypt_result, err := crypt.AesCBCDecrypt(decode_result)
	if err != nil {
		return nil, err
	}
	fmt.Println("response body:", string(decrypt_result))
	return decrypt_result, nil
}
