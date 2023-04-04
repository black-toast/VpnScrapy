package hulk

type Xpath struct {
	NovelTitle  string
	NovelAuthor string
	NovelDesc   string
}

var xpath *Xpath

func init() {
	xpath = new(Xpath)
	xpath.NovelTitle = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/h3"
	xpath.NovelAuthor = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/ul/li[1]/a"
	xpath.NovelDesc = "//*[@id=\"tab-description\"]/div"
}
