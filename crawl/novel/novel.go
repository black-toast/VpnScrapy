package novel

import (
	"VpnScrapy/crawl/novel/hulk"
	"VpnScrapy/http"
	"VpnScrapy/storage"
	"VpnScrapy/util"
	"fmt"
	"os"
	"strings"
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

func Scrapy(startChapter, endChapter int, path string) {
	if startChapter < -1 || endChapter < -1 {
		fmt.Println("start or end chapter is illegal chapter")
		return
	}

	saveNovelIntroductSwitch := false

	// request novel introduction
	fmt.Println("request novel introduction")
	novelUrl := BuildNovelUrl() + path
	content, err := Request(novelUrl, "GET")
	if err != nil {
		panic(err)
	}
	novelScrapy := new(hulk.HulkScrapy)
	doc := novelScrapy.CreateParseDoc(string(content))
	title := novelScrapy.ParseTitle(doc)
	novelId := novelScrapy.ParseNovelId(doc)

	titlePath := strings.ReplaceAll(strings.ReplaceAll(title, " ", ""), ":", "")
	novelDir := fmt.Sprintf("%s%s%s", url.saveIntroductionPath, string(os.PathSeparator), titlePath)

	// save novel introduction
	if saveNovelIntroductSwitch {
		author := novelScrapy.ParseAuthor(doc)
		desc := novelScrapy.ParseDesc(doc)
		novelScrapy.Save(novelDir, FileNameTitle, title)
		novelScrapy.Save(novelDir, FileNameAuthor, author)
		novelScrapy.Save(novelDir, FileNameDesc, desc)
	}

	// request novel chapter list
	fmt.Println("wait 5s, then request novel chapter list")
	time.Sleep(5 * time.Second)
	novelChapterUrl := fmt.Sprintf(UrlHulkNovelDomain+PathChapterList, novelId)
	content, err = Request(novelChapterUrl, "GET")
	if err != nil {
		panic(err)
	}
	doc = novelScrapy.CreateParseDoc(string(content))
	chapterList := novelScrapy.ParseNovelChapterList(doc)
	if chapterListJson, err := util.Format(chapterList); err != nil {
		panic(err)
	} else {
		novelScrapy.Save(novelDir, FileChapterList, chapterListJson)
	}

	if len(chapterList) <= 0 {
		fmt.Println("chapter list is empty.")
		return
	}
	fmt.Println("chapter list number is", len(chapterList))

	if startChapter <= -1 {
		startChapter = 1
	}

	if endChapter <= -1 {
		endChapter = len(chapterList)
	}

	// request novel chapter
	for index, chapter := range chapterList {
		if index+1 < startChapter {
			continue
		}
		if index+1 > endChapter {
			break
		}
		if index == 0 {
			fmt.Println("wait 2s, then request novel chapter(", chapter.Name, ")")
		} else {
			fmt.Println("wait 2s, then request next novel chapter(", chapter.Name, ")")
		}
		time.Sleep(2 * time.Second)

		t := time.Now()
		//参数必须是这个时间,格式任意
		// s := t.Format("2006-01-02 15:04:05")
		currentTime := t.Format("2006-01-02 15:04:05")
		fmt.Printf("current time: %s\n", currentTime)

		startCost := t.Unix()
		content, err = Request(chapter.Url, "GET")
		if err != nil {
			panic(err)
		}
		doc = novelScrapy.CreateParseDoc(string(content))
		chapterTitle := novelScrapy.ParseChapterTitle(doc)
		lineChapterTitle, chapter := novelScrapy.ParseChapter(doc)
		if chapterTitle == "" {
			chapterTitle = lineChapterTitle
		}

		// save chapter
		chapter = fmt.Sprintf("Chapter %d: %s\n%s", index+1, chapterTitle, chapter)
		saveChapterFileName := fmt.Sprintf(FileNameChapter, index+1)
		novelScrapy.Save(novelDir, saveChapterFileName, chapter)
		fmt.Println("save chapter complete")

		// transform mp3 audio
		ttsInput := fmt.Sprintf(`%s\ch-%d.txt`, novelDir, index+1)
		ttsMp3Output := fmt.Sprintf(`%s\ch-%d.mp3`, novelDir, index+1)
		util.EdgeTts(ttsInput, ttsMp3Output)
		fmt.Println("transform mp3 complete")

		// transform mp4 video
		videoImage := fmt.Sprintf(`%s\cover.jpg`, novelDir)
		ttsMp4Output := fmt.Sprintf(`%s\ch-%d.mp4`, novelDir, index+1)
		// 保证执行ffmpeg命令不会提示该文件已存在
		storage.Delete(ttsMp4Output)
		util.MakeImageVideo(videoImage, ttsMp3Output, ttsMp4Output)
		fmt.Println("transform mp4 complete")

		// storage.Delete(ttsInput)
		storage.Delete(ttsMp3Output)
		// storage.Delete(ttsMp4Output)

		endCost := time.Now().Unix()
		fmt.Printf("👆======================cost %ds=======================👆\n", (endCost - startCost))
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
