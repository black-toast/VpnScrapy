package hulk

import (
	"VpnScrapy/storage"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type HulkScrapy struct {
}

type ChapterList struct {
	Name string
	Url  string
}

func (scrapy HulkScrapy) CreateParseDoc(content string) *html.Node {
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		fmt.Println("parse doc fail")
		return nil
	}
	return doc
}

func (scrapy HulkScrapy) ParseTitle(doc *html.Node) string {
	node := htmlquery.FindOne(doc, NovelTitle)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseAuthor(doc *html.Node) string {
	node := htmlquery.FindOne(doc, NovelAuthor)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseDesc(doc *html.Node) string {
	node := htmlquery.FindOne(doc, NovelDesc)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseNovelId(doc *html.Node) string {
	node := htmlquery.FindOne(doc, NovelNovelId)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseReadChapterId(doc *html.Node, replaceStr string) string {
	node := htmlquery.FindOne(doc, NovelReadUrl)
	readUrl := strings.Trim(htmlquery.InnerText(node), "\n ")
	return strings.ReplaceAll(readUrl, replaceStr, "")
}

func (scrapy HulkScrapy) ParseNovelChapterList(doc *html.Node) []*ChapterList {
	urlNodes := htmlquery.Find(doc, NovelChapterListUrl)
	nameNodes := htmlquery.Find(doc, NovelChapterListName)

	urlLen := len(urlNodes)
	nameLen := len(nameNodes)
	var chapterList []*ChapterList
	if urlLen != nameLen {
		return chapterList
	}

	for index, nameNode := range nameNodes {
		name := strings.Trim(htmlquery.InnerText(nameNode), "\n ")
		url := strings.Trim(htmlquery.InnerText(urlNodes[index]), "\n ")
		chapterList = append(chapterList, &ChapterList{
			Name: name,
			Url:  url,
		})
	}
	return chapterList
}

func (scrapy HulkScrapy) ParseChapterTitle(doc *html.Node) string {
	node := htmlquery.FindOne(doc, NovelChapterTitle)
	return strings.Trim(htmlquery.InnerText(node), "\n ")
}

func (scrapy HulkScrapy) ParseChapter(doc *html.Node) string {
	nodes := htmlquery.Find(doc, NovelChapter)
	lines := len(nodes)
	chapter := ""
	for index, node := range nodes {
		line := removeSpecialChars(strings.Trim(htmlquery.InnerText(node), "\n "))
		if line == "" {
			continue
		}
		chapter += line
		if index+1 != lines {
			chapter += "\n"
		}
	}
	return chapter
}

func removeSpecialChars(str string) string {
	return strings.Replace(str, ".....", "", -1)
}

func (scrapy HulkScrapy) Save(fileDir, fileName, content string) {
	storage.Mkdir(fileDir)
	novelIntroduction := fmt.Sprintf("%s\\%s", fileDir, fileName)
	saveFile := storage.Create(novelIntroduction)
	saveFile.WriteString(content)
	saveFile.Close()
}
