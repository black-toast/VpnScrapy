package v1

import (
	"VpnScrapy/bean"
	vpnHttp "VpnScrapy/http"
	"VpnScrapy/util"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Launch(engine *gin.Engine) {
	v1Group := engine.Group("v1")
	{
		v1Group.GET("/vmess", generateVmessList)
		v1Group.GET("/ssr", generateSsList)
	}
}

func generateVmessList(c *gin.Context) {
	loginResp, err := vpnHttp.Login(
		util.WithConstant().GetTlyEmail(),
		util.WithConstant().GetTlyPassword(),
	)
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

	vmessList := ""
	for _, node := range result.Node {
		tlsConfig := ""
		if node.Tls == "true" {
			tlsConfig = "tls"
		}
		vmessEncodeSrc := fmt.Sprintf(
			util.WithConstant().GetVmessFormat(),
			node.NodeName,
			node.NodeServer,
			node.Port,
			node.Pass,
			node.NodeMethod,
			node.Cipher,
			node.WsPath,
			tlsConfig,
		)

		vmessList += fmt.Sprintf("vmess://%s\n", base64.StdEncoding.EncodeToString([]byte(vmessEncodeSrc)))
	}
	c.String(http.StatusOK, vmessList)
}

func generateSsList(c *gin.Context) {
	ssrSubscribeByte, err := vpnHttp.SsrSubscribe()
	if err != nil {
		fmt.Println("request ssr subscribe failure")
		return
	}

	result, err := util.Parse(ssrSubscribeByte, []bean.SsrSubscribe{})
	if err != nil {
		fmt.Println("parse ssr subscribe response data failure, err=", err)
		return
	}
	if len(result) <= 0 {
		fmt.Println("ssr vpn node list is empty")
		return
	}

	ssrSubscribeList := ""
	for _, node := range result {
		ssrLink := fmt.Sprintf(
			"%s:%d:origin:%s:plain:%s/?obfsparam=&remarks=%s&group=Y29tcGFueQ",
			node.Server,
			node.ServerPort,
			node.Method,
			encodeBase64(node.Password),
			encodeBase64(node.Remarks),
		)

		ssrSubscribeList += fmt.Sprintf("ssr://%s\n", encodeBase64(ssrLink))
	}
	c.String(http.StatusOK, ssrSubscribeList)
}

func encodeBase64(str string) string {
	encode := base64.StdEncoding.EncodeToString([]byte(str))
	return strings.ReplaceAll(encode, "=", "")
}
