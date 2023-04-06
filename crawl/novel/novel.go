package novel

import (
	"VpnScrapy/crawl/novel/hulk"
	"VpnScrapy/http"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/html"
)

type NovelScrapy interface {
	CreateParseDoc(content string) *html.Node
	ParseTitle(doc *html.Node) string
	ParseAuthor(doc *html.Node) string
	ParseDesc(doc *html.Node) string
	Save(fileDir, fileName, content string)
}

func Scrapy(startCrawlChapter int) {
	// request novel introduction
	fmt.Println("request novel introduction")
	novelUrl := BuildNovelUrl()
	content, err := Request(novelUrl, "GET")
	if err != nil {
		panic(err)
	}
	novelScrapy := new(hulk.HulkScrapy)
	doc := novelScrapy.CreateParseDoc(string(content))
	title := novelScrapy.ParseTitle(doc)
	author := novelScrapy.ParseAuthor(doc)
	desc := novelScrapy.ParseDesc(doc)
	novelId := novelScrapy.ParseNovelId(doc)
	readChapterId := novelScrapy.ParseReadChapterId(doc, novelUrl+"/")

	// save novel introduction
	novelDir := fmt.Sprintf("%s%s%s", url.saveIntroductionPath, string(os.PathSeparator), title)
	novelScrapy.Save(novelDir, FileNameTitle, title)
	novelScrapy.Save(novelDir, FileNameAuthor, author)
	novelScrapy.Save(novelDir, FileNameDesc, desc)

	// request novel chapter list
	fmt.Println("wait 5s, then request novel chapter list")
	time.Sleep(5 * time.Second)
	novelChapterUrl := fmt.Sprintf(UrlHulkNovelDomain+PathChapterList, novelId, readChapterId)
	content, err = Request(novelChapterUrl, "GET")
	if err != nil {
		panic(err)
	}
	doc = novelScrapy.CreateParseDoc(string(content))
	chapterList := novelScrapy.ParseNovelChapterList(doc)

	if len(chapterList) <= 0 {
		fmt.Println("chapter list is empty.")
		return
	}

	// request novel chapter
	for index, chapter := range chapterList {
		if index+1 < startCrawlChapter {
			continue
		}
		if index == 0 {
			fmt.Println("wait 5s, then request novel chapter(", chapter.Name, ")")
		} else {
			fmt.Println("wait 5s, then request next novel chapter(", chapter.Name, ")")
		}
		time.Sleep(5 * time.Second)

		content, err = Request(chapter.Url, "GET")
		if err != nil {
			panic(err)
		}
		doc = novelScrapy.CreateParseDoc(string(content))
		chapterTitle := novelScrapy.ParseChapterTitle(doc)
		chapter := novelScrapy.ParseChapter(doc)

		// chapterSplit := strings.Split(chapter, "\n")
		// line := chapterSplit[0]
		// lineLower := strings.ToLower(line)
		// chapterTitleLower := strings.ToLower(chapterTitle)
		// index := strings.Index(lineLower, chapterTitleLower+"1")
		// fmt.Println("lineLower:", lineLower, ",chapterTitleLower:", chapterTitleLower, ",index:", index)

		// fmt.Printf("%s\n%s\n%s", chapterTitle, nextChapterUrl, chapter)
		saveChapterFileName := fmt.Sprintf(FileNameChapter, index+1)
		chapter = chapterTitle + "\n" + chapter
		novelScrapy.Save(novelDir, saveChapterFileName, chapter)
	}
}

func Request(url string, method string) ([]byte, error) {
	novelContent, err := http.Request(&http.RequestConfig{
		Method:    method,
		Url:       url,
		Transport: http.V2rayProxy(),
	})
	if err != nil {
		fmt.Printf("request %s fail.\n %v\n", url, err)
		return nil, err
	}

	return novelContent, nil
}
