package novel

import (
	"fmt"
	"os"
)

const (
	FileNameTitle      = "title.txt"
	FileNameAuthor     = "author.txt"
	FileNameDesc       = "desc.txt"
	FileNameChapter    = "ch-%d.txt"
	UrlHulkNovelDomain = "https://novelhulk.com"
	PathNovelName      = "/nb/spy-mage-system-book"
	PathChapterList    = "/ajax/chapter-option?novelId=%s&currentChapterId=%s"
)

type Url struct {
	saveIntroductionPath string
}

var url *Url

func init() {
	url = new(Url)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	url.saveIntroductionPath = fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), "output")
}

func WithUrl() *Url {
	return url
}

func BuildNovelUrl() string {
	return UrlHulkNovelDomain + PathNovelName
}
