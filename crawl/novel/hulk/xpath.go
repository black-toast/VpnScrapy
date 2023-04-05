package hulk

const (
	// 小说详情页
	NovelTitle   = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/h3"
	NovelAuthor  = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/ul/li[1]/a"
	NovelDesc    = "//*[@id=\"tab-description\"]/div"
	NovelNovelId = "//div[@id=\"rating\"]/@data-novel-id"
	NovelReadUrl = "//a[contains(@class, \"btn-read-now\")]/@href"

	// 小说目录
	NovelChapterListUrl  = "//option/@value"
	NovelChapterListName = "//option/text()"

	// 小说章节页
	NovelChapterTitle = "//a[@class=\"chr-title\"]/span"
	NovelChapter      = "//*[@id=\"chr-content\"]/p"
)
