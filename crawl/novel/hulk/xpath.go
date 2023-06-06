package hulk

const (
	// 小说详情页
	NovelTitle   = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/h3"
	NovelAuthor  = "//div[contains(@class, \"desc\")]/ul/li[2]/a/text()"
	NovelDesc    = "//*[@id=\"tab-description\"]/div"
	NovelNovelId = "//div[@id=\"rating\"]/@data-novel-id"
	NovelReadUrl = "//a[contains(@class, \"btn-read-now\")]/@href"

	// 小说目录
	NovelChapterListUrl  = "//a/@href"
	NovelChapterListName = "//a/@title"

	// 小说章节页
	NovelChapterTitle = "//a[@class=\"chr-title\"]/span"
	NovelChapter      = "//*[@id=\"chr-content\"]/p"
)
