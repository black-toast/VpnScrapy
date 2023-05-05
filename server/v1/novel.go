package v1

import (
	"VpnScrapy/crawl/novel"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LaunchNovel(group *gin.RouterGroup) {
	novelGroup := group.Group("novel")
	{
		novelGroup.GET("/crawlNovel", crawlNovel)
	}
}

// 爬取小说
func crawlNovel(c *gin.Context) {
	sQuery := c.DefaultQuery("s", "0")
	startCrawlChapter, err := strconv.Atoi(sQuery)
	if err != nil {
		panic(err)
	}
	go novel.Scrapy(startCrawlChapter, -1)
}
