package novel

import (
	"VpnScrapy/crawl/novel/hulk"
	"VpnScrapy/storage"
	"fmt"
	"golang.org/x/net/html"
)

type NovelScrapy interface {
	CreateParseDoc(content string) *html.Node
	ParseTitle(doc *html.Node) string
	ParseAuthor(doc *html.Node) string
	ParseDesc(doc *html.Node) string
}

func Scrapy(content string) {
	novelScrapy := new(hulk.HulkScrapy)
	doc := novelScrapy.CreateParseDoc(content)
	title := novelScrapy.ParseTitle(doc)
	author := novelScrapy.ParseAuthor(doc)
	desc := novelScrapy.ParseDesc(doc)

	//fmt.Printf("title: %s author: %s\ndesc: %s", title, author, desc)
	novelDir := url.saveIntroductionPath + "\\" + title
	storage.Mkdir(novelDir)
	novelIntroduction := novelDir + "\\Introduction"
	saveFile := storage.Create(novelIntroduction)
	saveFile.WriteString(fmt.Sprintf("%s\n%s\n%s", title, author, desc))
	saveFile.Close()
}
