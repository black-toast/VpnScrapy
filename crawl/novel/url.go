package novel

import (
	"fmt"
	"os"
)

const (
	FileNameTitle      = "title.txt"
	FileNameAuthor     = "author.txt"
	FileNameDesc       = "desc.txt"
	FileChapterList    = "chapter_list.txt"
	FileNameChapter    = "ch-%d.txt"
	UrlHulkNovelDomain = "https://novelhulk.com"
	PathNovelName      = "/nb/necropolis-immortal-book"
	PathChapterList    = "/ajax/chapter-archive?novelId=%s"
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
	url.saveIntroductionPath = fmt.Sprintf("%s%s%s%s%s", wd, string(os.PathSeparator),
		"output", string(os.PathSeparator), "novels")
}

func WithUrl() *Url {
	return url
}

func BuildNovelUrl() string {
	// return UrlHulkNovelDomain + PathNovelName
	return UrlHulkNovelDomain + "/nb/"
}
