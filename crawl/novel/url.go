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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	url.saveIntroductionPath = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), "Desktop")
}

func WithUrl() *Url {
	return url
}

func BuildNovelUrl() string {
	return UrlHulkNovelDomain + PathNovelName
}
