package novel

import "os"

type Url struct {
	saveIntroductionPath string
	urlHulkNovelDomain   string
	pathNovelName        string
}

var url *Url

func init() {
	url = new(Url)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	url.saveIntroductionPath = homeDir + "\\Desktop"
	url.urlHulkNovelDomain = "https://novelhulk.com"
	url.pathNovelName = "/nb/spy-mage-system-book"
}

func WithUrl() *Url {
	return url
}

func (url *Url) BuildNovelUrl() string {
	return url.urlHulkNovelDomain + url.pathNovelName
}
