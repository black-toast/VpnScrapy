package v1

import (
	"VpnScrapy/bean"
	vpnHttp "VpnScrapy/http"
	"VpnScrapy/util"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Launch(engine *gin.Engine) {
	v1Group := engine.Group("v1")
	{
		v1Group.GET("/vmess", generateVmessList)
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
