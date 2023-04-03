package main

import (
	"VpnScrapy/bean"
	"VpnScrapy/http"
	"VpnScrapy/util"
	"encoding/base64"
	"fmt"
)

var vmessFormat = "{\"v\": \"2\",\"ps\": \"%s\",\"add\": \"%s\",\"port\": \"%d\",\"id\": \"%s\",\"aid\": \"0\",\"scy\": \"auto\",\"net\": \"%s\",\"type\": \"%s\",\"host\": \"\",\"path\": \"%s\",\"tls\": \"%s\",\"sni\": \"\",\"alpn\": \"\"}"

func main() {
	// tly
	// 1479405751@qq.com
	// ----1479405751@qq.com---pG5cOwEmHZXyMIMqBXOEHMAdAWrUD5QqVCqyL6s3Hg8=
	// 1479405751yy
	loginResp, err := http.Login("1479405751@qq.com", "1479405751yy")
	if err != nil {
		fmt.Println("exec login failure")
		return
	}

	result, err := util.Parse(loginResp, bean.TlyLoginResp{})
	if err != nil {
		fmt.Println("parse login response data failure, err=", err)
		return
	}
	if result.NodeNumber <= 0 || len(result.Node) <= 0 {
		fmt.Println("vpn node list is empty")
		return
	}

	for _, node := range result.Node {
		tlsConfig := ""
		if node.Tls == "true" {
			tlsConfig = "tls"
		}
		vmessEncodeSrc := fmt.Sprintf(
			vmessFormat,
			node.NodeName,
			node.NodeServer,
			node.Port,
			node.Pass,
			node.NodeMethod,
			node.Cipher,
			node.WsPath,
			tlsConfig,
		)

		fmt.Printf("vmess://%s\n", base64.StdEncoding.EncodeToString([]byte(vmessEncodeSrc)))
	}

}
