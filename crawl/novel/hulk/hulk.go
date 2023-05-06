package hulk

import (
	"VpnScrapy/storage"
	"errors"
	"fmt"
	"regexp"
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
	return parseRealChapterTitle(strings.Trim(htmlquery.InnerText(node), "\n "))
}

func (scrapy HulkScrapy) ParseChapter(doc *html.Node) (string, string) {
	nodes := htmlquery.Find(doc, NovelChapter)
	lines := len(nodes)
	chapter := ""
	parseChapterTitleLine := 0
	chapterTitle := ""
	for index, node := range nodes {
		lineContent := strings.Trim(htmlquery.InnerText(node), "\n ")

		// 解析第一行标题
		if index == parseChapterTitleLine {
			if lineContent == "" {
				parseChapterTitleLine++
				continue
			}
			chapterTitle = parseFirstLineChapterTitle(lineContent)
			if chapterTitle != "" {
				// fmt.Printf("find line chapter title: %s", chapterTitle)
				continue
			}
		}

		line := removeSpecialChars(lineContent)
		if isEndLine(line) {
			break
		}
		if line == "" {
			continue
		}
		chapter += line
		if index+1 != lines {
			chapter += "\n"
		}
	}
	return chapterTitle, chapter
}

// 解析真正的章节标题
// @return chapterNumber, chapterTitle
func parseRealChapterTitle(content string) string {
	// extract chapter number
	// example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-1
	compileRegex := regexp.MustCompile(`(\d{1,}).?-?:?`)
	matchArr := compileRegex.FindStringSubmatch(content)
	if len(matchArr) < 2 {
		return content
	}
	chapterNum := matchArr[len(matchArr)-1]

	chapterNumStartIndex := strings.LastIndex(content, chapterNum)
	if chapterNumStartIndex == -1 {
		panic(errors.New("illegal chapter title format(" + content + ")"))
	}

	chapterTitleStartIndex := chapterNumStartIndex + len(chapterNum)
	chapterTitle := strings.Trim(content[chapterTitleStartIndex:], ": ")
	chapterTitle = strings.Trim(chapterTitle, "- ")
	return chapterTitle
}

// 解析文章内容第一行中可能存在的章节标题
// @return chapterNumber, chapterTitle
func parseFirstLineChapterTitle(content string) string {
	// extract chapter number
	// example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-1
	compileRegex := regexp.MustCompile(`(\d{1,}).?-?:?`)
	matchArr := compileRegex.FindStringSubmatch(content)
	if len(matchArr) < 2 {
		return ""
	}
	chapterNum := matchArr[len(matchArr)-1]

	chapterNumStartIndex := strings.LastIndex(content, chapterNum)
	if chapterNumStartIndex == -1 {
		panic(errors.New("illegal line chapter title format(" + content + ")"))
	}

	chapterTitleStartIndex := chapterNumStartIndex + len(chapterNum)
	chapterTitle := strings.Trim(content[chapterTitleStartIndex:], ": ")
	chapterTitle = strings.Trim(chapterTitle, "- ")
	return chapterTitle
}

func removeSpecialChars(content string) string {
	// remove ..... format
	// example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-1
	compileRegex := regexp.MustCompile(`(\.{2,})`)
	matchArr := compileRegex.FindStringSubmatch(content)
	if len(matchArr) >= 2 {
		return ""
	}

	// example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-27
	if strings.HasPrefix(content, "The source of this content is") {
		return ""
	}

	if strings.HasPrefix(content, "Translator:") {
		return ""
	}

	compileRegex = regexp.MustCompile(`(…{2,})`)
	matchArr = compileRegex.FindStringSubmatch(content)
	if len(matchArr) >= 2 {
		return ""
	}

	compileRegex = regexp.MustCompile(`(\*{3,})`)
	matchArr = compileRegex.FindStringSubmatch(content)
	if len(matchArr) >= 2 {
		return ""
	}

	if content == "[…]" || content == "-" || content == "“…”" || content == "__" {
		return ""
	}

	webContent := strings.ReplaceAll(content, "[", "")
	webContent = strings.ReplaceAll(webContent, "]", "")
	webContent = strings.Trim(strings.ReplaceAll(webContent, "/", ""), " ")
	if strings.Contains(webContent, ".com") || strings.Contains(webContent, ".net") {
		return ""
	}

	return content
}

func isEndLine(content string) bool {
	if content == "Note:" || content == "Notes:" || content == "Endnote:" || content == "Endnote" {
		return true
	}
	if content == "Preview:" {
		return true
	}
	return false
}

func (scrapy HulkScrapy) Save(fileDir, fileName, content string) {
	storage.Mkdir(fileDir)
	novelIntroduction := fmt.Sprintf("%s\\%s", fileDir, fileName)
	saveFile := storage.Create(novelIntroduction)
	saveFile.WriteString(content)
	saveFile.Close()
}
