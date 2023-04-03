package server

import (
	v1 "VpnScrapy/server/v1"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConfigServer() {
	ginEngine := gin.Default()
	ginEngine.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Vpn Scrapy")
	})

	v1.Launch(ginEngine)

	// listen and serve on 0.0.0.0:8999 (for windows "localhost:8999")
	err := ginEngine.Run(":8999")
	if err != nil {
		fmt.Println("launch server failure")
		return
	}
}
