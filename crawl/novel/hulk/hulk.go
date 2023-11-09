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
	Name string `json:"name"`
	Url  string `json:"url"`
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
	return parseNovelDesc(htmlquery.InnerText(node))
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
	if lines < 10 {
		nodes = htmlquery.Find(doc, NovelChapter2)
		lines = len(nodes)
	}
	chapter := ""
	parseChapterTitleLine := 0
	chapterTitle := ""

chapterNodesFor:
	for index, node := range nodes {
		if len(node.Attr) > 0 {
			for _, attr := range node.Attr {
				// 跳过不展示的p标签内容
				if attr.Key == "style" && attr.Val == "display: none;" {
					continue chapterNodesFor
				}
			}
		}
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

// 解析小说描述
func parseNovelDesc(desc string) string {
	desc = strings.Trim(desc, "\n ")
	descSplit := strings.Split(desc, "\n")
	desc = ""
	for index, split := range descSplit {
		tmp := strings.ReplaceAll(split, "…", "")
		tmp = strings.ReplaceAll(tmp, "-", "")
		if tmp == "" || split == "" || strings.LastIndex(split, "Translator") == 0 {
			continue
		}
		if strings.LastIndex(split, "Disclaimer") == 0 ||
			strings.LastIndex(split, "Follow me on") == 0 ||
			strings.Contains(split, "@Webnovel_MLB") {
			continue
		}
		if strings.LastIndex(split, "Find out in") == 0 {
			break
		}

		if index != 0 {
			desc += "\n"
		}
		desc += split
	}
	return desc
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
	content = strings.ReplaceAll(content, `\"`, `"`)
	// remove ..... format
	// example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-1
	compileRegex := regexp.MustCompile(`^"?\.{2,}"?$`)
	matchArr := compileRegex.FindStringSubmatch(content)
	if len(matchArr) > 0 {
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

	if content == "[…]" || content == "-" || content == "“…”" || content == "__" ||
		content == "–" || content == "—" || content == "…..." || content == "…" {
		return ""
	}

	webContent := strings.ReplaceAll(content, "[", "")
	webContent = strings.ReplaceAll(webContent, "]", "")
	webContent = strings.Trim(strings.ReplaceAll(webContent, "/", ""), " ")
	compileRegex = regexp.MustCompile(`(\.[c|n|𝑪|𝓬|𝐂|𝑐|𝐜|𝕔|𝒸|𝗰|𝔠|𝚌|𝓒|𝒞][o|e|𝞸|𝚘|𝑶|𝓞|𝔬|𝐎|𝗈|𝒐|𝑂|𝑜|𝒪][m|t|𝓜|𝑚|𝚖|𝔪|𝓶|𝑀|𝗆|𝕞|𝓂|𝐦|𝐌])`)
	matchArr = compileRegex.FindStringSubmatch(webContent)
	if len(matchArr) >= 2 {
		return ""
	}

	return content
}

func isEndLine(content string) bool {
	if content == "Note:" || content == "Notes:" || strings.Contains(content, `Author’s Note:`) ||
		strings.Contains(content, "[Notes:") || strings.Contains(content, "notes:") {
		return true
	}
	if content == "Endnote:" || content == "Endnote" || content == "Endnotes:" {
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
