package bean

type Novel struct {
	Index       int             `json:"index"`
	Name        string          `json:"name"`
	ChapterList []*NovelChapter `json:"chapters"`
}

type NovelChapter struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	DESC  string `json:"desc"`
}
