package v1

import (
	"VpnScrapy/crawl/novel"
	"VpnScrapy/crawl/novel/hulk"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LaunchNovel(group *gin.RouterGroup) {
	novelGroup := group.Group("novel")
	{
		novelGroup.GET("/crawl_novel", crawlNovel)
		novelGroup.GET("/check_update", checkNovelsUpdate)
	}
}

// 检查
func checkNovelsUpdate(c *gin.Context) {
	novelScrapy := new(hulk.HulkScrapy)
	novelUpdateResult := "["
	checkNovelUpdate(novelScrapy, "spy-mage-system", &novelUpdateResult, true)
	checkNovelUpdate(novelScrapy, "necropolis-immortal", &novelUpdateResult, true)
	checkNovelUpdate(novelScrapy, "cultivation-online-novel", &novelUpdateResult, true)
	checkNovelUpdate(novelScrapy, "the-experimental-log-of-the-crazy-lich", &novelUpdateResult, true)
	checkNovelUpdate(novelScrapy, "invincible-divine-dragons-cultivation-system", &novelUpdateResult, true)
	checkNovelUpdate(novelScrapy, "trial-marriage-husband-need-to-work-hard", &novelUpdateResult, false)
	novelUpdateResult += "]"
	c.String(http.StatusOK, novelUpdateResult)
}

func checkNovelUpdate(hulkScrapy *hulk.HulkScrapy, novelId string, result *string, appendComma bool) {
	novelChapterUrl := fmt.Sprintf(novel.UrlHulkNovelDomain+novel.PathChapterList, novelId)
	content, err := novel.Request(novelChapterUrl, "GET")
	if err != nil {
		panic(err)
	}
	doc := hulkScrapy.CreateParseDoc(string(content))
	chapterList := hulkScrapy.ParseNovelChapterList(doc)

	if len(chapterList) != 0 {
		chapterSize := len(chapterList)
		*result += fmt.Sprintf(`{"novelId": "%s", "index": %d, "lastest": "%s"}`, novelId, chapterSize, chapterList[chapterSize-1].Name)
	} else {
		*result += fmt.Sprintf(`{"novelId": "%s", "lastest": "none"}`, novelId)
	}
	if appendComma {
		*result += ","
	}
}

// 爬取小说
func crawlNovel(c *gin.Context) {
	sQuery := c.DefaultQuery("s", "0")
	startCrawlChapter, err := strconv.Atoi(sQuery)
	if err != nil {
		panic(err)
	}
	go novel.Scrapy(startCrawlChapter, -1, "")
}
